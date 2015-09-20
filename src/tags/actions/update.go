package tagactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// HandleUpdateShow serves a get request at /tags/1/update (show form to update)
func HandleUpdateShow(context router.Context) error {

	// Find the tag
	tag, err := tags.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, tag)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("tag", tag)
	return view.Render()
}

// HandleUpdate serves POST or PUT /tags/1/update
func HandleUpdate(context router.Context) error {

	// Find the tag
	tag, err := tags.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, tag)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the tag
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	err = tag.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to tag
	return router.Redirect(context, tag.URLShow())
}
