package imageactions

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// POST /images/1/destroy
func HandleDestroy(context router.Context) error {

	// Set the image on the context for checking
	image, err := images.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, image)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the image
	image.Destroy()

	// Redirect to sites - better than images root as a default
	// but should never be used if we have a redirect
	return router.Redirect(context, "/sites")

}
