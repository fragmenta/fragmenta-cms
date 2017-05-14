package pageactions

import (
	"net/http"

	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/pages"
	"github.com/fragmenta/fragmenta-cms/src/users"
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
	page, err := pages.FindFirst("url=?", "/")
	if err != nil {
		return server.NotFoundError(nil)
	}

	currentUser := session.CurrentUser(w, r)

	view := view.NewWithPath(r.URL.Path, w)
	view.AddKey("title", "Fragmenta app")
	view.AddKey("page", page)
	view.AddKey("currentUser", currentUser)
	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.Template("pages/views/templates/default.html.got")
	return view.Render()
}
