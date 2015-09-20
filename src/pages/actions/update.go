package pageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// HandleUpdateShow handles GET /pages/1/update (show form to update)
func HandleUpdateShow(context router.Context) error {

	// Find the page
	page, err := pages.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error Error finding page %s", err)
		return router.NotFoundError(err)
	}

	// Authorise updating page
	err = authorise.Resource(context, page)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("page", page)
	return view.Render()
}

// HandleUpdate handles POST or PUT /pages/1/update
func HandleUpdate(context router.Context) error {

	// Find the page
	page, err := pages.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise updating the page
	err = authorise.Resource(context, page)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the page from params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}
	err = page.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// We then find the page again, and retreive the new Url, in case it has changed during update
	page, err = pages.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Redirect to page url
	return router.Redirect(context, page.Url)
}
