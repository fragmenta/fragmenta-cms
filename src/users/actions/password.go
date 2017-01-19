package useractions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/mail"
	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

const (
	// ResetLifetime is the maximum time reset tokens are valid for
	ResetLifetime = time.Hour
)

// HandlePasswordResetShow responds to GET /users/password/reset
// by showing the password reset page.
func HandlePasswordResetShow(w http.ResponseWriter, r *http.Request) error {
	// No authorisation required, just show the view
	view := view.NewRenderer(w, r)
	view.Template("users/views/password_reset.html.got")
	return view.Render()
}

// HandlePasswordResetSend responds to POST /users/password/reset
// by sending a password reset email.
func HandlePasswordResetSend(w http.ResponseWriter, r *http.Request) error {

	// No authorisation required
	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return server.NotAuthorizedError(err, "Invalid authenticity token")
	}

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user by email (if not found let them know)
	// Find the user by hex token in the db
	email := params.Get("email")
	user, err := users.FindFirst("email=?", email)
	if err != nil {
		return server.Redirect(w, r, "/users/password/reset?message=invalid_email")
	}

	// Generate a random token and url for the email
	token := auth.BytesToHex(auth.RandomToken(32))

	// Update the user record with with this token
	userParams := map[string]string{
		"password_reset_token": token,
		"password_reset_at":    query.TimeString(time.Now().UTC()),
	}
	// Direct access to the user columns, bypassing validation
	user.Update(userParams)

	// Generate the url to use in our email
	url := fmt.Sprintf("%s/users/password?token=%s", config.Get("root_url"), token)

	// Send a password reset email out to this user
	emailContext := map[string]interface{}{
		"url":  url,
		"name": user.Name,
	}

	log.Info(log.V{"msg": "sending reset email", "user_email": user.Email, "user_id": user.ID})

	e := mail.New(user.Email)
	e.Subject = "Reset Password"
	e.Template = "users/views/password_reset_mail.html.got"
	err = mail.Send(e, emailContext)
	if err != nil {
		return err
	}

	// Tell the user what we have done
	return server.Redirect(w, r, "/users/password/sent")
}

// HandlePasswordResetSentShow responds to GET /users/password/sent
func HandlePasswordResetSentShow(w http.ResponseWriter, r *http.Request) error {
	view := view.NewRenderer(w, r)
	view.Template("users/views/password_sent.html.got")
	return view.Render()
}

// HandlePasswordReset responds to GET /users/password?token=DEADFISH
// by logging the user in, removing the token
// and allowing them to set their password.
func HandlePasswordReset(w http.ResponseWriter, r *http.Request) error {

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Note we have no authenticity check, just a random token to check
	token := params.Get("token")
	if len(token) < 10 || len(token) > 64 {
		return server.NotAuthorizedError(fmt.Errorf("Invalid reset token"), "Invalid Token")
	}

	// Find the user by hex token in the db
	user, err := users.FindFirst("password_reset_token=?", token)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Make sure the reset at time is less expire time
	if time.Since(user.PasswordResetAt) > ResetLifetime {
		return server.NotAuthorizedError(nil, "Token invalid", "Your password reset token has expired, please request another.")
	}

	// Remove the reset token from this user
	// using direct access, bypassing validation
	user.Update(map[string]string{"password_reset_token": ""})

	// Log in the user and store in the session
	// Now save the user details in a secure cookie, so that we remember the next request
	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(w, r)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Save login the session cookie
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(w)

	// Log action
	log.Info(log.V{"msg": "reset password", "user_email": user.Email, "user_id": user.ID})

	// Redirect to the user update page so that they can change their password
	return server.Redirect(w, r, fmt.Sprintf("/users/%d/update", user.ID))
}
