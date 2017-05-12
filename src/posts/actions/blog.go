package postactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/session"
	"github.com/fragmenta/fragmenta-cms/src/posts"
)

// HandleShowBlog responds to GET /blog
func HandleShowBlog(w http.ResponseWriter, r *http.Request) error {

	// Authorise list post
	user := session.CurrentUser(w, r)
	err := can.List(posts.New(), user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Build a query for blog posts in chronological order
	q := posts.Published().Order("created_at desc").Limit(50)
	blogPosts, err := posts.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	log.Log(log.V{"MSG": blogPosts})

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("posts", blogPosts)
	view.Template("posts/views/blog.html.got")
	return view.Render()
}
