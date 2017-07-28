package pageactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/pages"
	"github.com/fragmenta/fragmenta-cms/src/redirects"
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
	view.AddKey("meta_title", page.Name)
	view.AddKey("meta_keywords", page.Keywords)
	view.AddKey("meta_desc", page.Summary)
	view.Template(page.ShowTemplate())
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
	path := "/" + params.Get("path")
	page, err := pages.FindFirst("url=?", path)
	if err != nil {
		redirect, err := redirects.FindFirst("old_url=?", path)
		if err != nil {
			return server.NotFoundError(err)
		}
		return server.Redirect(w, r, redirect.NewURL)
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
	view.AddKey("meta_title", page.Name)
	view.AddKey("meta_keywords", page.Keywords)
	view.AddKey("meta_desc", page.Summary)
	view.Template(page.ShowTemplate())
	return view.Render()
}
