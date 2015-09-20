package pageactions

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// POST /pages/1/destroy
func HandleDestroy(context router.Context) error {

	// Set the page on the context for checking
	page, err := pages.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, page)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the page
	page.Destroy()

	// Redirect to pages root
	return router.Redirect(context, page.URLIndex())
}
