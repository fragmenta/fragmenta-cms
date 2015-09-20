package postactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// HandleShow displays a single post
func HandleShow(context router.Context) error {

	// Find the post
	post, err := posts.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Authorise access only if not published
	if !post.IsPublished() {
		err = authorise.Resource(context, post)
		if err != nil {
			return router.NotAuthorizedError(err)
		}
	}

	// Render the template
	view := view.New(context)
	view.AddKey("post", post)
	return view.Render()
}
