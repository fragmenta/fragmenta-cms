package app

import (
	"net/http"

	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"
)

// HandleShowHome serves our home page with a simple template.
// This function might be moved over to src/pages if you have a pages resource.
func homeHandler(w http.ResponseWriter, r *http.Request) error {

	// Demonstrate tracing in log messages
	log.Info(log.Values{"msg": "Home handler", "trace": log.Trace(r)})

	view := view.NewWithPath(r.URL.Path, w)
	view.AddKey("title", "Fragmenta app")
	view.Template("app/views/home.html.got")
	return view.Render()
}

// simpleHomeHandler demonstrates a simple handler without using the view pkg
// all packages in fragmenta are optional, if you want to use a different
// template library to render, don't import github.com/fragmenta/view.
func simpleHomeHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)
	w.Write([]byte("hello world"))
	return nil
}
