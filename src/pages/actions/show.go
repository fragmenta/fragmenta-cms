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

// HandleShow displays a single page.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise access
	user := session.CurrentUser(w, r)
	err = can.Show(page, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(page.CacheKey())
	view.AddKey("page", page)
	view.AddKey("currentUser", user)
	view.Template(page.Template)
	return view.Render()
}

// HandleShowPath serves requests to a custom page url
func HandleShowPath(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the page
	page, err := pages.FindFirst("url=?", "/"+params.Get("path"))
	if err != nil {

		// If no pages or users exist, redirect to set up page
		if missingUsersAndPages() {
			return server.Redirect(w, r, "/fragmenta/setup")
		}

		return server.NotFoundError(err)
	}

	// Authorise access IF the page is not published
	user := session.CurrentUser(w, r)
	if !page.IsPublished() {
		err = can.Show(page, user)
		if err != nil {
			return server.NotAuthorizedError(err)
		}
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(page.CacheKey())
	view.AddKey("page", page)
	view.AddKey("currentUser", user)
	//	view.Template(page.Template)
	view.Template("pages/views/templates/show.html.got")
	return view.Render()
}

func missingUsersAndPages() bool {

	// Check if there are zero pages and zero users.
	return false
}
