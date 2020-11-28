package mux

import (
	"fmt"
	"io"
	"net/http"
)

// fileHandler is the default static file handler called if there is no route.
func fileHandler(w http.ResponseWriter, r *http.Request) error {
	// Just return a not found error
	// Set the headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	// Write a simple error message page - omit error details for security reasons
	html := fmt.Sprintf("<h1>404 Not Found Error</h1>")
	io.WriteString(w, html)
	return nil
}

// errHandler is a simple built-in error handler which writes the error string to context.Writer
// users of the mux should override this with their own handler.
func errHandler(w http.ResponseWriter, r *http.Request, err error) {

	// Set the headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	// Write a simple error message page - omit error details for security reasons
	html := fmt.Sprintf("<h1>500 Internal Error</h1>")
	io.WriteString(w, html)
}
