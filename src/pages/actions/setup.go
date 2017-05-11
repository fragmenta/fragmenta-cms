package pageactions

import (
	"net/http"

	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/daintydomains/src/users"
)

// HandleShowSetup responds to GET /fragmenta/setup
func HandleShowSetup(w http.ResponseWriter, r *http.Request) error {

	// Check we have no users, if not bail out
	if users.Count() != 0 {
		return server.NotAuthorizedError(nil, "Users already exist")
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.Template("pages/views/setup.html.got")
	return view.Render()
}

// HandleSetup responds to POST /fragmenta/setup
func HandleSetup(w http.ResponseWriter, r *http.Request) error {

	// Check we have no users, if not bail out
	if users.Count() != 0 {
		return server.NotAuthorizedError(nil, "Users already exist")
	}

	// Set up a new user using the form details

	// Set up a new home page

	/*

	   <section class="padded">
	   <h1>Welcome to Fragmenta</h1>
	   <h2>Edit this page</h2>

	   <a href="/pages/1/update">Edit this page</a>
	   </section>

	*/

	// Redirect to the home page (newly set up we hope)
	return server.Redirect(w, r, "/")
}
