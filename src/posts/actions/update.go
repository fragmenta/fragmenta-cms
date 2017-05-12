package postactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/posts"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleUpdateShow renders the form to update a post.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the post
	post, err := posts.Find(params.GetInt(posts.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise update post
	user := session.CurrentUser(w, r)
	err = can.Update(post, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Fetch the users
	authors, err := users.FindAll(users.Where("role=?", users.Admin))
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("post", post)
	view.AddKey("authors", authors)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a post
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the post
	post, err := posts.Find(params.GetInt(posts.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update post
	user := session.CurrentUser(w, r)
	err = can.Update(post, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	postParams := post.ValidateParams(params.Map(), posts.AllowedParams())

	err = post.Update(postParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to post
	return server.Redirect(w, r, post.ShowURL())
}
