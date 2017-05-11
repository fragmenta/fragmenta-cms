package pageactions

import (
	"net/http"

	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/kennygrant/daintydomains/src/pages"
	"github.com/kennygrant/daintydomains/src/users"
)

// HandleShowHome serves our home page with a simple template.
func HandleShowHome(w http.ResponseWriter, r *http.Request) error {

	// Demonstrate tracing in log messages
	log.Info(log.Values{"msg": "Home handler", "trace": log.Trace(r)})

	// If we have no users (first run), redirect to setup
	if users.Count() == 0 {
		return server.Redirect(w, r, "/fragmenta/setup")
	}

	// Home fetches the first page with the url '/' and uses it for the home page of the site
	page, err := pages.FindFirst("path='/'")
	if err != nil {
		return server.NotFoundError(nil)
	}

	view := view.NewWithPath(r.URL.Path, w)
	view.AddKey("title", "Fragmenta app")
	view.AddKey("page", page)
	view.Template("pages/views/home.html.got")
	return view.Render()
}
