package postactions

import (
	"net/http"

	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// HandleShowBlog responds to GET /blog
func HandleShowBlog(w http.ResponseWriter, r *http.Request) error {

	// Build a query for blog posts in chronological order
	q := posts.Published().Order("created_at desc").Limit(50)
	blogPosts, err := posts.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	user := session.CurrentUser(w, r)

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("posts", blogPosts)
	view.Template("posts/views/blog.html.got")
	return view.Render()
}
