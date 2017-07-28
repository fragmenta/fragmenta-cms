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

// HandleCreateShow serves the create form via GET for redirects.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	redirect := redirects.New()

	// Authorise
	user := session.CurrentUser(w, r)
	err := can.Create(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("redirect", redirect)
	return view.Render()
}

// HandleCreate handles the POST of the create form for redirects
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	redirect := redirects.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	user := session.CurrentUser(w, r)
	err = can.Create(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	redirectParams := redirect.ValidateParams(params.Map(), redirects.AllowedParams())

	id, err := redirect.Create(redirectParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new redirect
	redirect, err = redirects.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, redirect.IndexURL())
}
