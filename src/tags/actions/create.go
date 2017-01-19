package tagactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// HandleCreateShow serves the create form via GET for tags.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	tag := tags.New()

	// Authorise
	err := can.Create(tag, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("tag", tag)
	return view.Render()
}

// HandleCreate handles the POST of the create form for tags
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	tag := tags.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	err = can.Create(tag, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	tagParams := tag.ValidateParams(params.Map(), tags.AllowedParams())

	id, err := tag.Create(tagParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new tag
	tag, err = tags.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, tag.IndexURL())
}
