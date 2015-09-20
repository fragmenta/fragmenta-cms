package postactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// HandleBlog displays a list of posts in reverse chronological order
func HandleBlog(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Build a query
	q := posts.Published().Order("created_at desc")

	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("( to_tsvector(name) || to_tsvector(summary) @@ to_tsquery(?) )", filter)
	}

	// Fetch the posts
	results, err := posts.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	context.Logf("POSTS HERE :%v", results)

	// Render the template
	view := view.New(context)
	view.AddKey("filter", filter)
	view.AddKey("posts", results)
	view.Template("posts/views/blog.html.got")
	return view.Render()

}
