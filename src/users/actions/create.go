package useractions

import (
	"fmt"
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleCreateShow serves the create form via GET for users.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	user := users.New()

	fmt.Printf("USER:%v\n", user)

	// Authorise
	currentUser := session.CurrentUser(w, r)
	err := can.Create(user, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", currentUser)
	view.AddKey("user", user)
	return view.Render()
}

// HandleCreate handles the POST of the create form for users
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	user := users.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	currentUser := session.CurrentUser(w, r)
	err = can.Create(user, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	userParams := user.ValidateParams(params.Map(), users.AllowedParams())

	id, err := user.Create(userParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new user
	user, err = users.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, user.IndexURL())
}
