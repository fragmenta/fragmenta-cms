package pageactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/pages"
)

// HandleIndex serves a get request at /pages
func HandleIndex(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Fetch the pages
	q := pages.Query().Order("url asc")

	// Filter if necessary
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("(to_tsvector(name) || to_tsvector(summary) || to_tsvector(url) @@ to_tsquery(?) )", filter)
	}

	pageList, err := pages.FindAll(q)
	if err != nil {
		context.Logf("#error Error indexing pages %s", err)
		return router.InternalError(err)
	}

	// Serve template
	view := view.New(context)
	view.AddKey("filter", filter)
	view.AddKey("pages", pageList)
	return view.Render()

}
