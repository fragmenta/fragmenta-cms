package postactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/posts"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleUpdateShow renders the form to update a post
func HandleUpdateShow(context router.Context) error {

	// Find the post
	post, err := posts.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update post
	err = authorise.Resource(context, post)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Find users for author details
	users, err := users.FindAll(users.Admins())
	if err != nil {
		return router.NotFoundError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("post", post)
	view.AddKey("users", users)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a post
func HandleUpdate(context router.Context) error {

	// Find the post
	post, err := posts.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update post
	err = authorise.Resource(context, post)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the post from params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}
	err = post.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to post
	return router.Redirect(context, post.URLShow())
}
