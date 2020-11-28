package log

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RequestID is but a simple token for tracing requests.
type RequestID struct {
	id []byte
}

// String returns a string formatting for the request id.
func (r *RequestID) String() string {
	return fmt.Sprintf("%X-%X-%X-%X", r.id[0:2], r.id[2:4], r.id[4:6], r.id[6:8])
}

// NewRequestID returns a new random request id.
func newRequestID() *RequestID {
	r := &RequestID{
		id: make([]byte, 8),
	}
	rand.Read(r.id)
	return r
}

type ctxKey struct{}

// Trace retreives the request id from a request as a string.
func Trace(r *http.Request) string {
	rid, ok := r.Context().Value(&ctxKey{}).(*RequestID)
	if ok {
		return rid.String()
	}
	return ""
}

// GetRequestID retreives the request id from a request.
func GetRequestID(r *http.Request) *RequestID {
	return r.Context().Value(&ctxKey{}).(*RequestID)
}

// SetRequestID saves the request id in the request context.
func SetRequestID(r *http.Request, rid *RequestID) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, &ctxKey{}, rid)
	return r.WithContext(ctx)
}

// Middleware adds a logging wrapper and request tracing to requests.
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		r = SetRequestID(r, requestID) // Sets on context for handlers

		level := LevelInfo

		// For assets etc, use level debug as they clutter up logs
		if r.URL.Path == "/favicon.ico" ||
			strings.HasPrefix(r.URL.Path, "/assets") ||
			strings.HasPrefix(r.URL.Path, "/stats") {
			level = LevelDebug
		}

		Log(Values{
			MessageKey: "<- Request",
			"method":   r.Method,
			URLKey:     r.RequestURI,
			"len":      r.ContentLength,
			IPKey:      r.RemoteAddr,
			TraceKey:   requestID.String(),
			LevelKey:   level,
		})

		start := time.Now()
		h(w, r)

		Time(start, Values{
			MessageKey: "-> Response",
			URLKey:     r.RequestURI,
			TraceKey:   requestID.String(),
			LevelKey:   level,
		})
	}

}
