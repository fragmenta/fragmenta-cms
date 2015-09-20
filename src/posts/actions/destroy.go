package postactions

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// HandleDestroy handles a DESTROY request for posts
func HandleDestroy(context router.Context) error {

	// Find the post
	post, err := posts.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise destroy post
	err = authorise.Resource(context, post)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the post
	post.Destroy()

	// Redirect to posts root
	return router.Redirect(context, post.URLIndex())
}
