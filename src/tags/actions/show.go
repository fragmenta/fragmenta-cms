package tagactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// HandleShow displays a single tag.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the tag
	tag, err := tags.Find(params.GetInt(tags.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	err = can.Show(tag, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(tag.CacheKey())
	view.AddKey("tag", tag)
	return view.Render()
}
