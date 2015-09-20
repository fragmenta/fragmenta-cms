package pageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// GET pages/create
func HandleCreateShow(context router.Context) error {
	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup
	view := view.New(context)
	page := pages.New()
	view.AddKey("page", page)

	// Serve
	return view.Render()
}

// POST pages/create
func HandleCreate(context router.Context) error {
	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	id, err := pages.Create(params.Map())

	if err != nil {
		context.Logf("#info Failed to create page %v", params)
		return router.InternalError(err)
	}

	// Log creation
	context.Logf("#info Created page id,%d", id)

	// Redirect to the new page
	p, err := pages.Find(id)
	if err != nil {
		context.Logf("#error Error creating page,%s", err)
	}

	return router.Redirect(context, p.URLIndex())
}
