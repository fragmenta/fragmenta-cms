package pageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// HandleShow shows pages using their ID: GET /pages/1
func HandleShow(context router.Context) error {

	// Find the page
	page, err := pages.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Authorise access
	err = authorise.Resource(context, page)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the page
	return renderPage(context, page)
}

// HandleShowPath serves requests to a custom page url
func HandleShowPath(context router.Context) error {

	// Setup context for template
	path := context.Path()

	// If no pages or users exist, redirect to set up page
	if missingUsersAndPages() {
		return router.Redirect(context, "/fragmenta/setup")
	}

	q := pages.Query().Where("url=?", path).Limit(1)
	pages, err := pages.FindAll(q)
	if err != nil || len(pages) == 0 {
		return router.NotFoundError(err)
	}

	// Get the first of pages to render
	page := pages[0]

	// For show path of pages, we authorise showing the page FOR ALL users if it is published
	if !page.IsPublished() {
		// Authorise
		err = authorise.Resource(context, page)
		if err != nil {
			return router.NotAuthorizedError(err)
		}
	}

	return renderPage(context, page)
}

func renderPage(context router.Context, page *pages.Page) error {

	view := view.New(context)

	// Setup context for template
	if page.Template != "" {
		view.Template(page.Template)
	} else {
		view.Template("pages/views/show.html.got")
	}

	view.AddKey("page", page)
	view.AddKey("meta_title", page.Name)
	view.AddKey("meta_desc", page.Summary)
	view.AddKey("meta_keywords", page.Keywords)

	// Serve template
	context.Logf("#info Rendering page for path %s", context.Path())
	return view.Render()
}
