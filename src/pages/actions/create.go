package pageactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// HandleCreateShow serves the create form via GET for pages.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	page := pages.New()

	// Authorise
	err := can.Create(page, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("page", page)
	return view.Render()
}

// HandleCreate handles the POST of the create form for pages
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	page := pages.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	err = can.Create(page, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	pageParams := page.ValidateParams(params.Map(), pages.AllowedParams())

	id, err := page.Create(pageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new page
	page, err = pages.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, page.IndexURL())
}
