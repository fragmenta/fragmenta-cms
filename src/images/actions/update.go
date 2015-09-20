package imageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// Serve a get request at /images/1/update (show form to update)
func HandleUpdateShow(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	image, err := images.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error Error finding image %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, image)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	view.AddKey("redirect", context.Param("redirect"))
	view.AddKey("image", image)
	view.AddKey("admin_links", helpers.Link("Destroy Image", url.Destroy(image), "method=delete"))

	return view.Render()
}

// POST or PUT /images/1/update
func HandleUpdate(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	// Find the image
	image, err := images.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error Error finding image %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, image)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the image
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	err = image.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// We redirect back to source if redirect param is set
	return router.Redirect(context, url.Update(image))

}
