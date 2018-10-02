package server

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// StatusError wraps a std error and stores more information (status code, display title/msg and caller info)
type StatusError struct {
	Err     error
	Status  int
	Title   string
	Message string
	File    string
	Line    int
}

// Error returns the underling error string - it should not be shown in production
func (e *StatusError) Error() string {
	return fmt.Sprintf("Status %d at %s : %s", e.Status, e.FileLine(), e.Err)
}

// String returns a string represenation of this error, useful for debugging
func (e *StatusError) String() string {
	return fmt.Sprintf("Status %d at %s : %s %s %s", e.Status, e.FileLine(), e.Title, e.Message, e.Err)
}

// FileLine returns file name and line of error
func (e *StatusError) FileLine() string {
	parts := strings.Split(e.File, "/")
	f := strings.Join(parts[len(parts)-4:len(parts)], "/")
	return fmt.Sprintf("%s:%d", f, e.Line)
}

func (e *StatusError) setupFromArgs(args ...string) *StatusError {
	if e.Err == nil {
		e.Err = fmt.Errorf("Error:%d", e.Status)
	}
	if len(args) > 0 {
		e.Title = args[0]
	}
	if len(args) > 1 {
		e.Message = args[1]
	}
	return e
}

// NotFoundError returns a new StatusError with Status StatusNotFound and optional Title and Message
// Usage return router.NotFoundError(err,"Optional Title", "Optional user-friendly Message")
func NotFoundError(e error, args ...string) *StatusError {
	err := Error(e, http.StatusNotFound, "Not Found", "Sorry, the page you're looking for couldn't be found.")
	return err.setupFromArgs(args...)
}

// InternalError returns a new StatusError with Status StatusInternalServerError and optional Title and Message
// Usage: return router.InternalError(err)
func InternalError(e error, args ...string) *StatusError {
	err := Error(e, http.StatusInternalServerError, "Server Error", "Sorry, something went wrong, please let us know.")
	return err.setupFromArgs(args...)
}

// NotAuthorizedError returns a new StatusError with Status StatusUnauthorized and optional Title and Message
func NotAuthorizedError(e error, args ...string) *StatusError {
	err := Error(e, http.StatusUnauthorized, "Not Allowed", "Sorry, I can't let you do that.")
	return err.setupFromArgs(args...)
}

// BadRequestError returns a new StatusError with Status StatusBadRequest and optional Title and Message
func BadRequestError(e error, args ...string) *StatusError {
	err := Error(e, http.StatusBadRequest, "Bad Request", "Sorry, there was an error processing your request, please check your data.")
	return err.setupFromArgs(args...)
}

// Error returns a new StatusError with code StatusInternalServerError and a generic message
func Error(e error, s int, t string, m string) *StatusError {
	// Get runtime info - use zero values if none available
	_, f, l, _ := runtime.Caller(2)
	err := &StatusError{
		Status:  s,
		Err:     e,
		Title:   t,
		Message: m,
		File:    f,
		Line:    l,
	}
	return err
}

// ToStatusError returns a *StatusError or wraps a standard error in a 500 StatusError
func ToStatusError(e error) *StatusError {
	if err, ok := e.(*StatusError); ok {
		return err
	}
	return Error(e, http.StatusInternalServerError, "Error", "Sorry, an error occurred.")
}
