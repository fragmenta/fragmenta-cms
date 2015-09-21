package imageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// HandleShow serves a get request at /images/1
func HandleShow(context router.Context) error {

	// Find the image
	image, err := images.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Authorise
	err = authorise.Resource(context, image)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Serve template
	view := view.New(context)
	view.AddKey("image", image)
	return view.Render()

}
