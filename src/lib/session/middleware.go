package session

import (
	"context"
	"net/http"
	"strings"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"
)

// Middleware sets a token on every GET request so that it can be
// inserted into the view. It currently ignores requests for files and assets.
func Middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// If a get method, we need to set the token for use in views
		if shouldSetToken(r) {
			// This sets the token on the encrypted session cookie
			// and the request Context (for use in view later?)
			token, err := auth.AuthenticityToken(w, r)
			if err != nil {
				log.Error(log.Values{"msg": "session: problem setting token", "error": err})
			} else {
				// Save the token to the request context for use in views
				ctx := r.Context()
				ctx = context.WithValue(ctx, view.AuthenticityContext, token)
				r = r.WithContext(ctx)
			}

		}

		h(w, r)
	}

}

// shouldSetToken returns true if this request requires a token set.
func shouldSetToken(r *http.Request) bool {

	// No tokens on anything but GET requests
	if r.Method != http.MethodGet {
		return false
	}

	// No tokens on non-html resources
	if strings.HasPrefix(r.URL.Path, "/files") ||
		strings.HasPrefix(r.URL.Path, "/assets") {
		return false
	}

	return true
}
