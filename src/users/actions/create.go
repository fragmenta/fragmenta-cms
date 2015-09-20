package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// GET users/create
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup
	view := view.New(context)
	user := users.New()
	view.AddKey("user", user)

	// Serve
	return view.Render()
}

// POST users/create
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

	// We should check for duplicates in here
	id, err := users.Create(params.Map())
	if err != nil {
		return router.InternalError(err, "Error", "Sorry, an error occurred creating the user record.")
	} else {
		context.Logf("#info Created user id,%d", id)
	}

	// Redirect to the new user
	p, err := users.Find(id)
	if err != nil {
		return router.InternalError(err, "Error", "Sorry, an error occurred finding the new user record.")
	}

	return router.Redirect(context, p.URLIndex())
}
