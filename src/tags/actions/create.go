package tagactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// GET tags/create
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup
	view := view.New(context)
	tag := tags.New()
	view.AddKey("tag", tag)

	// Serve
	return view.Render()
}

// POST tags/create
func HandleCreate(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	id, err := tags.Create(params.Map())
	if err != nil {
		context.Logf("#info Failed to create tag %v", params)
		return router.InternalError(err)
	}

	// Log creation
	context.Logf("#info Created tag id,%d", id)

	// Redirect to the new tag
	tag, err := tags.Find(id)
	if err != nil {
		context.Logf("#error Error creating tag,%s", err)
	}

	// Always regenerate dotted ids - we fetch all tags first to avoid db calls
	q := tags.Query().Select("select id,parent_id from tags").Order("id asc")
	tagsList, err := tags.FindAll(q)
	if err == nil {
		dotted_params := map[string]string{}
		dotted_params["dotted_ids"] = tag.CalculateDottedIds(tagsList)
		tags.Query().Where("id=?", tag.Id).Update(dotted_params)
	}

	return router.Redirect(context, tag.URLIndex())
}
