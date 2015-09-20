package imageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// GET images/create
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup
	view := view.New(context)
	image := images.New()
	view.AddKey("image", image)

	// Serve
	return view.Render()
}

// POST images/create
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

	delete(params, "tab")

	id, err := images.Create(params.Map())
	if err != nil {
		context.Logf("#error Error creating image,%s", err)
		return router.InternalError(err)
	}

	// Log creation
	context.Logf("#info Created image id,%d", id)

	// Redirect to the new image
	m, err := images.Find(id)
	if err != nil {
		context.Logf("#error Error creating image,%s", err)
	}

	return router.Redirect(context, m.URLShow())
}
