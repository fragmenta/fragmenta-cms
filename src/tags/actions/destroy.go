package tagactions

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// POST /tags/1/destroy
func HandleDestroy(context router.Context) error {

	// Set the tag on the context for checking
	tag, err := tags.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, tag)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the tag
	tag.Destroy()

	// Redirect to tags root
	return router.Redirect(context, tag.URLIndex())
}
