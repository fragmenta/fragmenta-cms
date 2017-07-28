package redirectactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/redirects"
)

// HandleDestroy responds to /redirects/n/destroy by deleting the redirect.
func HandleDestroy(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the redirect
	redirect, err := redirects.Find(params.GetInt(redirects.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise destroy redirect
	user := session.CurrentUser(w, r)
	err = can.Destroy(redirect, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Destroy the redirect
	redirect.Destroy()

	// Redirect to redirects root
	return server.Redirect(w, r, redirect.IndexURL())

}
