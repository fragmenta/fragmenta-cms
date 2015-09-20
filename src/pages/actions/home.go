package pageactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
	"github.com/fragmenta/fragmenta-cms/src/users"
)

// HandleHome serves a get request at /
func HandleHome(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// If nothing exists, redirect to set up page
	if missingUsersAndPages() {
		return router.Redirect(context, "/fragmenta/setup")
	}

	page, err := pages.Find(1)
	if err != nil {
		return router.InternalError(err)
	}

	view.AddKey("page", page)
	view.AddKey("meta_title", page.Name)
	view.AddKey("meta_desc", page.Summary)
	view.AddKey("meta_keywords", page.Keywords)

	return view.Render()

}

func resultSet(a []*users.User) []int64 {
	var results []int64

	for _, u := range a {
		results = append(results, u.ImageID)
	}

	return results
}
