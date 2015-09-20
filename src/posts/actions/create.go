package postactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/posts"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleCreateShow serves the create form via GET for posts
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
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
	post := posts.New()

	user := authorise.CurrentUser(context)
	if user != nil {
		post.AuthorId = user.Id
	}

	view.AddKey("post", post)
	view.AddKey("users", users)
	return view.Render()
}

// HandleCreate handles the POST of the create form for posts
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

	id, err := posts.Create(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Log creation
	context.Logf("#info Created post id,%d", id)

	// Redirect to the new post
	m, err := posts.Find(id)
	if err != nil {
		return router.InternalError(err)
	}

	return router.Redirect(context, m.URLIndex())
}
