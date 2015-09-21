package tagactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// HandleShow serves a get request at /tags/1
func HandleShow(context router.Context) error {

	// Find the resource
	id := context.ParamInt("id")
	tag, err := tags.Find(id)
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, tag)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Serve template
	view := view.New(context)
	view.AddKey("tag", tag)
	return view.Render()
}
