package pageactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/pages"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleUpdateShow renders the form to update a page.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the page
	page, err := pages.Find(params.GetInt(pages.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise update page
	err = can.Update(page, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Fetch the users
	authors, err := users.FindAll(users.Where("role=?", users.Admin))
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("page", page)
	view.AddKey("authors", authors)
	view.AddKey("currentUser", session.CurrentUser(w, r))
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a page
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the page
	page, err := pages.Find(params.GetInt(pages.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update page
	err = can.Update(page, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	pageParams := page.ValidateParams(params.Map(), pages.AllowedParams())

	err = page.Update(pageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to page
	return server.Redirect(w, r, page.ShowURL())
}
