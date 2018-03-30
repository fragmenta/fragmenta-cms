package useractions

import (
	"fmt"
	"net/http"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleLoginShow shows the page at /users/login
func HandleLoginShow(w http.ResponseWriter, r *http.Request) error {

	// Check they're not logged in already.
	currentUser := session.CurrentUser(w, r)
	if !currentUser.Anon() {
		return server.Redirect(w, r, "/?warn=already_logged_in")
	}

	params, err := mux.Params(r)
	if err != nil {
		return server.NotFoundError(err)
	}

	// Show the login page, with login failure warnings.
	view := view.NewRenderer(w, r)
	switch params.Get("error") {
	case "failed_email":
		view.AddKey("warning", "Sorry, we couldn't find a user with that email.")
	case "failed_password":
		view.AddKey("warning", "Sorry, the password was incorrect, please try again.")
	}
	return view.Render()
}

// HandleLogin responds to POST /users/login
// by setting a cookie on the request with encrypted user data.
func HandleLogin(w http.ResponseWriter, r *http.Request) error {

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Check they're not logged in already if so redirect.
	currentUser := session.CurrentUser(w, r)
	if !currentUser.Anon() {
		return server.Redirect(w, r, "/?warn=already_logged_in")
	}

	// Get the user details from the database
	params, err := mux.Params(r)
	if err != nil {
		return server.NotFoundError(err)
	}

	email := params.Get("email")
	password := params.Get("password")

	// Fetch the first user by email
	user, err := users.FindFirst("email=?", email)
	if err != nil {
		log.Info(log.V{"msg": "login failed", "email": email, "status": http.StatusNotFound})
		return server.Redirect(w, r, "/users/login?error=failed_email")
	}

	// Check password against the stored password
	err = auth.CheckPassword(password, user.PasswordHash)
	if err != nil {
		log.Info(log.V{"msg": "login failed", "email": email, "user_id": user.ID, "status": http.StatusUnauthorized})
		return server.Redirect(w, r, "/users/login?error=failed_password")
	}

	// Now save the user details in a secure cookie, so that we remember the next request
	session, err := auth.Session(w, r)
	if err != nil {
		log.Info(log.V{"msg": "login failed", "email": email, "user_id": user.ID, "status": http.StatusInternalServerError})
	}

	// Success, log it and set the cookie with user id
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(w)

	// Log action
	log.Info(log.V{"msg": "login", "user_email": user.Email, "user_id": user.ID})

	// Redirect - ideally here we'd redirect to their original request path
	return server.Redirect(w, r, "/")
}
