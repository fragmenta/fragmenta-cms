package useractions

import (
	"fmt"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleUpdateShow serves a get request at /users/1/update (show form to update)
func HandleUpdateShow(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error Error finding user %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, user)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	view.AddKey("user", user)
	//	view.AddKey("admin_links", helpers.Link("Destroy User", url.Destroy(user), "method=post"))

	return view.Render()
}

// HandleUpdate or PUT /users/1/update
func HandleUpdate(context router.Context) error {

	// Find the user
	id := context.ParamInt("id")
	user, err := users.Find(id)
	if err != nil {
		context.Logf("#error Error finding user %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, user)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// We expect only one image, what about replacing the existing when updating?
	// At present we just create a new image
	files, err := context.ParamFiles("image")
	if err != nil {
		return router.InternalError(err)
	}

	// Get the params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	var imageID int64

	if len(files) > 0 {
		fileHandle := files[0]

		// Create an image (saving the image representation on disk)
		imageParams := map[string]string{"name": user.Name, "status": "100"}
		imageID, err = images.Create(imageParams, fileHandle)
		if err != nil {
			return router.InternalError(err)
		}

		params.Set("image_id", fmt.Sprintf("%d", imageID))
		delete(params, "image")
	}

	err = user.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to user
	return router.Redirect(context, user.URLShow())
}
