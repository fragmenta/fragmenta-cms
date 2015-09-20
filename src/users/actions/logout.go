package useractions

import (
	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"
)

// HandleLogout clears the current user's session /users/logout
func HandleLogout(context router.Context) error {

	// Clear the current session cookie
	auth.ClearSession(context)

	// Redirect to home
	return router.Redirect(context, "/")
}
