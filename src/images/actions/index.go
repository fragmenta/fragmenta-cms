package imageactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/images"
	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
)

// Serve a get request at /images
//
//
func HandleIndex(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context for template
	view := view.New(context)

	// Build a query
	q := images.Query()

	// Show only published (defaults to showing all)
	// q.Apply(status.WherePublished)

	// Order by required order, or default to id asc
	switch context.Param("order") {

	case "1":
		q.Order("created_at desc")

	case "2":
		q.Order("updated_at desc")

	case "3":
		q.Order("name asc")

	default:
		q.Order("id desc")

	}

	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("( to_tsvector(name) || to_tsvector(summary) @@ to_tsquery(?) )", filter)
	}

	// Fetch the images
	imagesList, err := images.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Serve template
	view.AddKey("filter", filter)
	view.AddKey("images", imagesList)
	// Can we add these programatically?
	view.AddKey("admin_links", helpers.Link("Create image", images.New().URLCreate()))

	return view.Render()

}
