package imageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// Serve a get request at /images/1

func HandleShow(context router.Context) error {

	// Setup context for template
	view := view.New(context)

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
	view.AddKey("image", image)
	view.AddKey("admin_links", helpers.Link("Edit Image", url.Update(image)))

	return view.Render()

}
