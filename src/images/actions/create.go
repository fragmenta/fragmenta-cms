package imageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// HandleCreateShow serves GET images/create
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Serve the template
	view := view.New(context)
	image := images.New()
	view.AddKey("image", image)
	return view.Render()
}

// HandleCreate responds to POST images/create
func HandleCreate(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// We expect only one image, what about replacing the existing when updating?
	// At present we just create a new image
	files, err := context.ParamFiles("image")
	if err != nil {
		return router.InternalError(err)
	}

	if len(files) == 0 {
		return router.NotFoundError(nil)
	}

	// Get the params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	id, err := images.Create(params.Map(), files[0])
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
