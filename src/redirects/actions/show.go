package redirectactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/redirects"
)

// HandleShow displays a single redirect.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the redirect
	redirect, err := redirects.Find(params.GetInt(redirects.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	user := session.CurrentUser(w, r)
	err = can.Show(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.CacheKey(redirect.CacheKey())
	view.AddKey("redirect", redirect)
	return view.Render()
}
