package useractions

import (
	"fmt"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/users"
)

// GET /users/password/reset
func HandlePasswordResetShow(context router.Context) error {
	view := view.New(context)
	view.Template("users/views/password_reset.html.got")
	return view.Render()
}

// POST /users/password/reset
func HandlePasswordResetSend(context router.Context) error {

	// Find the user by email (if not found let them know)
	// Find the user by hex token in the db
	email := context.Param("email")
	user, err := users.First(users.Where("email=?", email))
	if err != nil {
		return router.Redirect(context, "/users/password/reset?message=invalid_email")
	}

	// Generate a random token and url for the email
	token := auth.BytesToHex(auth.RandomToken())

	// Generate the url to use in our email
	base := fmt.Sprintf("%s://%s", context.Request().URL.Scheme, context.Request().URL.Host)
	url := fmt.Sprintf("%s/users/password?token=%s", base, token)

	context.Logf("#info sending reset email:%s url:%s", email, url)

	// Update the user record with with this token
	p := map[string]string{"reset_token": token}
	user.Update(p)

	// Send a password reset email out
	//mail.Send("mymail")

	// Tell the user what we have done
	return router.Redirect(context, "/users/password/sent")
}

// GET /users/password/sent
func HandlePasswordResetSentShow(context router.Context) error {
	view := view.New(context)
	view.Template("users/views/password.html.got")
	return view.Render()
}

// POST /users/password?token=DEADFISH - handle password reset link
func HandlePasswordReset(context router.Context) error {

	token := context.Param("token")
	if len(token) == 0 {
		return router.InternalError(fmt.Errorf("Blank reset token"))
	}

	// Find the user by hex token in the db
	user, err := users.First(users.Where("token=?", token))
	if err != nil {
		return router.InternalError(err)
	}

	// If we found a user with this token, log the sender in as them
	// and remove the token from the user so that it can't be used twice
	// we should possibly add a time limit to tokens too?

	// Redirect to the user update page so that they can change their password
	return router.Redirect(context, fmt.Sprintf("/users/%d/update", user.Id))
}
