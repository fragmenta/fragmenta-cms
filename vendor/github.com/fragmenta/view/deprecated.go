package view

import (
	"net/http"
)

// RenderContext is the type passed in to New, which helps construct the rendering view
// Alternatively, you can use NewWithPath, which doesn't require a RenderContext
type RenderContext interface {
	Path() string
	RenderContext() map[string]interface{}
	Writer() http.ResponseWriter
}

// New creates a new Renderer
func New(c RenderContext) *Renderer {
	r := &Renderer{
		path:     c.Path(),
		layout:   "app/views/layout.html.got",
		template: "",
		format:   "text/html",
		status:   http.StatusOK,
		context:  c.RenderContext(),
		writer:   c.Writer(),
	}

	// This sets layout and template based on the view.path
	r.setDefaultTemplates()

	return r
}

// NewWithPath creates a new Renderer with a path and an http.ResponseWriter
func NewWithPath(p string, w http.ResponseWriter) *Renderer {
	r := &Renderer{
		path:     p,
		layout:   "app/views/layout.html.got",
		template: "",
		format:   "text/html",
		status:   http.StatusOK,
		context:  make(map[string]interface{}, 0),
		writer:   w,
	}

	// This sets layout and template based on the view.path
	r.setDefaultTemplates()

	return r
}
