package mux

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// NewRoute returns a new Route of our default type.
func NewRoute(pattern string, handler HandlerFunc) (Route, error) {
	return NewPrefixRoute(pattern, handler)
}

// NewNaiveRoute creates a new Route, given a pattern to match and a handler for the route.
func NewNaiveRoute(pattern string, handler HandlerFunc) (Route, error) {
	r := &NaiveRoute{}
	err := r.Setup(pattern, handler)
	return r, err
}

// NewPrefixRoute creates a new PrefixRoute, given a pattern to match and a handler for the route.
func NewPrefixRoute(pattern string, handler HandlerFunc) (Route, error) {
	r := &PrefixRoute{}
	err := r.Setup(pattern, handler)
	return r, err
}

// NaiveRoute holds a pattern which matches a route and params within it,
// and an associated handler which will be called when the route matches.
type NaiveRoute struct {
	pattern    string
	handler    HandlerFunc
	methods    []string
	paramNames []string
	regexp     *regexp.Regexp
}

// Handler returns our handlerfunc.
func (r *NaiveRoute) Handler() HandlerFunc {
	return r.handler
}

// Setup sets up the route from a pattern
func (r *NaiveRoute) Setup(p string, h HandlerFunc) error {
	// Allow GET and HEAD by default
	r.methods = []string{http.MethodGet, http.MethodHead}
	r.handler = h
	r.pattern = p

	// Parse regexp once on startup
	return r.compileRegexp()
}

// Handle calls the handler with the writer and request.
func (r *NaiveRoute) Handle(w http.ResponseWriter, req *http.Request) error {
	return r.handler(w, req)
}

// MatchMethod returns true if our list of methods contains method
func (r *NaiveRoute) MatchMethod(method string) bool {

	for _, v := range r.methods {
		if v == method {
			return true
		}
		// Treat "" as GET
		if method == "" && v == http.MethodGet {
			return true
		}
	}

	return false
}

// MatchMaybe returns false if the path definitely is not MatchMethod
// or true/maybe if it *may* match.
func (r *NaiveRoute) MatchMaybe(path string) bool {
	return r.Match(path) // Just cheat and do a full match on base class
}

// Match returns true if this route matches the path given.
func (r *NaiveRoute) Match(path string) bool {

	// If we have a short pattern match, and we have a regexp, check against that
	if r.regexp != nil {
		return r.regexp.MatchString(path)
	}

	// If no regexp, check for exact string match against pattern
	return (r.pattern == path)
}

// Get sets the method exclusively to GET
func (r *NaiveRoute) Get() Route {
	return r.Method(http.MethodGet)
}

// Post sets the method exclusively to POST
func (r *NaiveRoute) Post() Route {
	return r.Method(http.MethodPost)
}

// Put sets the method exclusively to PUT
func (r *NaiveRoute) Put() Route {
	return r.Method(http.MethodPut)
}

// Delete sets the method exclusively to DELETE
func (r *NaiveRoute) Delete() Route {
	return r.Method(http.MethodDelete)
}

// Method sets the method exclusively to method
func (r *NaiveRoute) Method(method string) Route {
	r.methods = []string{method}
	return r
}

// Methods sets the methods allowed as an array
func (r *NaiveRoute) Methods(permitted ...string) Route {
	r.methods = permitted
	return r
}

// Pattern returns the string pattern for the route
func (r *NaiveRoute) Pattern() string {
	return r.pattern
}

// String returns the route formatted as a string
func (r *NaiveRoute) String() string {
	return fmt.Sprintf("%s %s", r.methods[0], r.pattern)
}

// Parse parses this path given our regexp and returns a map of URL params.
func (r *NaiveRoute) Parse(path string) map[string]string {

	// Set up our params map
	params := make(map[string]string, 0)

	// If called on a nil route, return empty params
	if r == nil || r.regexp == nil || len(r.paramNames) == 0 {
		return params
	}

	// Find a set of matches, and for each match set the entry in our map.
	matches := r.regexp.FindStringSubmatch(path)

	if matches != nil {
		for i, key := range r.paramNames {
			index := i + 1
			if len(matches) > index {
				value := matches[index]
				params[key] = value
			}
		}
	}

	return params
}

// compileRegexp compiles our route format to a true regexp
// Both name and regexp are required - routes should be well structured and restrictive by default
// Convert the pattern from the form  /pages/{id:[0-9]*}/edit
// to one suitable for regexp -  /pages/([0-9]*)/edit
// We want to match things like this:
// /pages/{id:[0-9]*}/edit
// /pages/{id:[0-9]*}/edit?param=test
func (r *NaiveRoute) compileRegexp() (err error) {

	// First return if no regexp
	if !strings.Contains(r.pattern, "{") {
		return nil
	}

	// Check if it is well-formed.
	idxs, errBraces := r.findBraces(r.pattern)
	if errBraces != nil {
		return errBraces
	}

	pattern := bytes.NewBufferString("^")
	end := 0

	// Walk through indexes two at a time
	for i := 0; i < len(idxs); i += 2 {
		// Set all values we are interested in.
		raw := r.pattern[end:idxs[i]]
		end = idxs[i+1]
		parts := strings.SplitN(r.pattern[idxs[i]+1:end-1], ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Missing name or pattern in %s", raw)
		}

		// Add the name to params in order of finding
		r.paramNames = append(r.paramNames, parts[0])

		// Add the real regexp
		fmt.Fprintf(pattern, "%s(%s)", regexp.QuoteMeta(raw), parts[1])

	}
	// Add the remaining pattern
	pattern.WriteString(regexp.QuoteMeta(r.pattern[end:]))
	r.regexp, err = regexp.Compile(pattern.String())

	return err
}

// findBraces returns the first level curly brace indices from a string.
// It returns an error in case of unbalanced braces.
// This method of parsing regexp is based on gorilla mux.
func (r *NaiveRoute) findBraces(s string) ([]int, error) {
	var level, idx int
	var idxs []int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("Route error: unbalanced braces in %q", s)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("Route error: unbalanced braces in %q", s)
	}
	return idxs, nil
}

// PrefixRoute uses a static prefix to reject route matches quickly.
type PrefixRoute struct {
	NaiveRoute
	index int
}

// Setup sets up the pattern prefix for the Prefix route.
func (r *PrefixRoute) Setup(p string, h HandlerFunc) error {

	// Record the prefix len up to the first regexp (if any)
	r.index = strings.Index(p, "{")

	// Finish setup with NaiveRoute
	return r.NaiveRoute.Setup(p, h)
}

// MatchMaybe returns false if the path definitely is not MatchMethod
// or true/maybe if it *may* match.
func (r *PrefixRoute) MatchMaybe(path string) bool {

	// If no prefix we are static, so can safely match absolutely
	if r.index < 0 {
		return path == r.pattern
	}

	// Reject with a string comparison of static prefix with path
	// HasPrefix checks on length first so it is fast.
	// If this returns yes, we are really saying maybe
	// and require a further check with Match().
	return strings.HasPrefix(path, r.pattern[:r.index])
}

// String returns the route formatted as a string.
func (r *PrefixRoute) String() string {
	return fmt.Sprintf("%s %s (prefix:%s)", r.methods[0], r.pattern, r.pattern[:r.index])
}
