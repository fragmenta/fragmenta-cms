package tagactions

import (
	"fmt"
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-cms/src/lib/authorise"
	"github.com/fragmenta/fragmenta-cms/src/tags"
)

// Only required once
func UpdateAllDottedIds() error {

	// Regenerate dotted ids - we fetch all tags first to avoid db calls
	q := tags.Query().Select("select id,parent_id from tags").Order("id asc")
	tagsList, err := tags.FindAll(q)
	if err == nil {
		for _, tag := range tagsList {
			params := map[string]string{
				"dotted_ids": tag.CalculateDottedIds(tagsList),
			}
			fmt.Printf("\n%d -> %s\n", tag.Id, params["dotted_ids"])
			err = tags.Query().Where("id=?", tag.Id).Update(params)
		}

	} else {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}

// HandleIndex serves a get request at /tags
func HandleIndex(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	//UpdateAllDottedIds()

	// Setup context for template
	view := view.New(context)

	// Fetch the tags
	q := tags.RootTags().Order("name asc")

	// Filter if necessary
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("(to_tsvector(name) || to_tsvector(summary) @@ to_tsquery(?))", filter)
	}

	tagList, err := tags.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Serve template
	view.AddKey("filter", filter)
	view.AddKey("tags", tagList)

	return view.Render()

}
