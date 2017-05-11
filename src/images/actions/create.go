package imageactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/session"
)

// HandleCreateShow serves the create form via GET for images.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	image := images.New()

	// Authorise
	user := session.CurrentUser(w, r)
	err := can.Create(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("image", image)
	return view.Render()
}

// HandleCreate handles the POST of the create form for images
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	image := images.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	user := session.CurrentUser(w, r)
	err = can.Create(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	imageParams := image.ValidateParams(params.Map(), images.AllowedParams())

	id, err := image.Create(imageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new image
	image, err = images.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, image.IndexURL())
}
