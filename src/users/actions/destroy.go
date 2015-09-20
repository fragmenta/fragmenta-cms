package useractions

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// POST /users/1/destroy
func HandleDestroy(context router.Context) error {

	// Set the user on the context for checking
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, user)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the user
	user.Destroy()

	// Redirect to users root
	return router.Redirect(context, user.URLIndex())
}
