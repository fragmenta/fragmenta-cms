package useractions

import (
	"fmt"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleShow serve a get request at /users/1
func HandleShow(context router.Context) error {

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error parsing user id: %s", err)
		return router.NotFoundError(err)
	}

	userMeta := fmt.Sprintf("%s – %s", user.Name, user.Summary)

	// Set up view
	view := view.New(context)

	// Find the first image which matches this user
	image, err := images.Find(user.ImageID)
	if err == nil {
		// only add image key if we have one
		view.AddKey("image", image)
	}

	// Render the Template
	view.AddKey("user", user)
	view.AddKey("meta_title", userMeta)
	view.AddKey("meta_desc", userMeta)
	view.AddKey("meta_keywords", user.Keywords())
	return view.Render()

}
