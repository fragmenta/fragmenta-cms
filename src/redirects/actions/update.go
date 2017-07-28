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

// HandleUpdateShow renders the form to update a redirect.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise update redirect
	user := session.CurrentUser(w, r)
	err = can.Update(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("redirect", redirect)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a redirect
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

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

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update redirect
	user := session.CurrentUser(w, r)
	err = can.Update(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	redirectParams := redirect.ValidateParams(params.Map(), redirects.AllowedParams())

	err = redirect.Update(redirectParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to redirect
	return server.Redirect(w, r, redirect.ShowURL())
}
