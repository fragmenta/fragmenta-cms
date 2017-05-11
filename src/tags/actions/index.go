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

// HandleIndex displays a list of tags.
func HandleIndex(w http.ResponseWriter, r *http.Request) error {

	// Authorise list tag
	user := session.CurrentUser(w, r)
	err := can.List(tags.New(), user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Build a query
	q := tags.Query()

	// Order by required order, or default to id asc
	switch params.Get("order") {

	case "1":
		q.Order("created desc")

	case "2":
		q.Order("updated desc")

	default:
		q.Order("id asc")
	}

	// Filter if requested
	filter := params.Get("filter")
	if len(filter) > 0 {
		q.Where("name ILIKE ?", filter)
	}

	// Fetch the tags
	results, err := tags.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("filter", filter)
	view.AddKey("tags", results)
	return view.Render()
}
