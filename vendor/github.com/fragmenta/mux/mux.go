package mux

import (
	"net/http"
	"strings"
	"sync"
)

// HandlerFunc defines a std net/http HandlerFunc, but which returns an error.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandlerFunc defines a HandlerFunc which accepts an error and displays it.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// Middleware is a handler that wraps another handler
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route defines the interface routes are expected to conform to.
type Route interface {
	// Match against URL
	MatchMethod(string) bool
	MatchMaybe(string) bool
	Match(string) bool

	// Handler returns the handler to execute
	Handler() HandlerFunc

	// Parse the URL for params according to pattern
	Parse(string) map[string]string

	// Set accepted methods
	Get() Route
	Post() Route
	Put() Route
	Delete() Route
	Methods(...string) Route
}

// MaxCacheEntries defines the maximum number of entries in the request->route cache
// 0 means caching is turned off
var MaxCacheEntries = 500

// mux is a private variable which is set only once on startup.
var mux *Mux

// SetDefault sets the default mux on the package for use in parsing params
// we could instead decorate each request with a reference to the Route
// but this means extra allocations for each request,
// when almost all apps require only one mux.
func SetDefault(m *Mux) {
	if mux == nil {
		mux = m

		// Set our router to handle all routes
		http.Handle("/", mux)
	}
}

// Mux handles http requests by selecting a handler
// and passing the request to it.
// Routes are evaluated in the order they were added.
// Before the request reaches the handler
// it is passed through the middleware chain.
type Mux struct {
	cache   map[string]Route
	cacheMu sync.RWMutex

	routes       []Route
	handlerFuncs []Middleware

	// See httptrace for best way to instrument
	ErrorHandler ErrorHandlerFunc
	FileHandler  HandlerFunc
	RedirectWWW  bool
}

// New returns a new mux
func New() *Mux {
	m := &Mux{
		RedirectWWW:  false,
		FileHandler:  fileHandler,
		ErrorHandler: errHandler,
		cache:        make(map[string]Route, MaxCacheEntries),
	}

	return m
}

// ServeHTTP implements net/http.Handler.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If redirect www is set, test host
	if m.RedirectWWW && strings.HasPrefix(r.Host, "www.") {
		redirect := strings.Replace("https://"+r.Host+r.URL.String(), "www.", "", 1)
		http.Redirect(w, r, redirect, http.StatusMovedPermanently)
	}

	// Avoid iteration if possible
	if len(m.handlerFuncs) == 0 {
		m.RouteRequest(w, r)
		return
	}
	h := m.RouteRequest
	for _, mh := range m.handlerFuncs {
		h = mh(h)
	}
	h(w, r)
}

// RouteRequest is the final endpoint of all requests
func (m *Mux) RouteRequest(w http.ResponseWriter, r *http.Request) {
	// Match a route
	route := m.Match(r)
	if route == nil {
		err := m.FileHandler(w, r)
		if err != nil {
			m.ErrorHandler(w, r, err)
		}
		return
	}

	// Execute the route
	err := route.Handler()(w, r)
	if err != nil {
		m.ErrorHandler(w, r, err)
	}

}

// Match finds the route (if any) which matches this request
func (m *Mux) Match(r *http.Request) Route {
	// Handle nil request
	if r == nil {
		return nil
	}

	// Check if we have a cached result for this same path
	if MaxCacheEntries > 0 {
		m.cacheMu.RLock()
		route, ok := m.cache[requestCacheKey(r)]
		m.cacheMu.RUnlock()
		// This check is necessary as we only use the request url for the cache key
		// this means we get a cache miss on some identical routes with diff methods.
		if ok && route.MatchMethod(r.Method) {
			return route
		}
	}

	// Routes are checked in order against the request path
	for _, route := range m.routes {
		// Test with probabalistic match
		if route.MatchMaybe(r.URL.Path) {
			// Test on method
			if route.MatchMethod(r.Method) {
				// Test exact match (may be expensive regexp)
				if route.Match(r.URL.Path) {
					m.cacheRoute(requestCacheKey(r), route)
					return route
				}
			}

		}
	}

	return nil
}

// Return a key suitable for storing this request in our cache.
// NB: To avoid allocations we do not include every permutation in the cache
// so routes returned must be checked against request.
func requestCacheKey(r *http.Request) string {
	return r.URL.Path
}

// cacheRoute saves the route with key provided
func (m *Mux) cacheRoute(key string, r Route) {
	if MaxCacheEntries == 0 {
		return // MaxCacheEntries is 0 so cache is off
	}
	m.cacheMu.Lock()
	// If cache is too big, evict
	if len(m.cache) > MaxCacheEntries {
		m.cache = make(map[string]Route, MaxCacheEntries)
	}
	// Fill the cache for this key -> route pair
	m.cache[key] = r
	m.cacheMu.Unlock()
}

// AddMiddleware adds a middleware function, this should be done before
// starting the server as it remakes our chain of middleware.
// This prepends to our chain of middleware
func (m *Mux) AddMiddleware(middleware Middleware) {
	m.handlerFuncs = append([]Middleware{middleware}, m.handlerFuncs...)
}

// AddHandler adds a route for this pattern using a
// stdlib http.HandlerFunc which does not return an error.
func (m *Mux) AddHandler(pattern string, handler http.HandlerFunc) Route {
	return m.Add(pattern, func(w http.ResponseWriter, r *http.Request) error {
		handler(w, r)
		return nil
	})
}

// Add adds a route for this request with the default methods (GET/HEAD)
// Route is returned so that method functions can be chained
func (m *Mux) Add(pattern string, handler HandlerFunc) Route {
	route, err := NewRoute(pattern, handler)
	if err != nil {
		// errors should be rare, but log them to stdout for debug
		println("mux: error parsing route:%s", pattern)
	}

	m.routes = append(m.routes, route)
	return route
}

// Get adds a route for this pattern/hanlder with the default methods (GET/HEAD)
func (m *Mux) Get(pattern string, handler HandlerFunc) Route {
	return m.Add(pattern, handler)
}

// Post adds a route for this pattern/hanlder with method http.PostMethod
func (m *Mux) Post(pattern string, handler HandlerFunc) Route {
	return m.Add(pattern, handler).Post()
}
