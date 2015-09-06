package app

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// HandleShowHome serves the home page - in a real app this might be elsewhere
func HandleShowHome(context router.Context) {
	view := view.New(context)

	view.AddKey("meta_title", "Fragmenta")
	view.AddKey("meta_desc", "Fragmenta App")
	view.AddKey("meta_keywords", "fragmenta, website")
	view.AddKey("title", "Hello world!")
	view.Template("app/views/home.html.got")
	view.Render(context)
}
