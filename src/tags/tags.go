// The tags package
package tags

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// The tags model type
type Tag struct {
	model.Model
	status.ModelStatus

	Name         string
	ParentId     int64
	Summary      string
	Url          string
	DisplayOrder int64
	DottedIds    string
}

// Which parameters does this model allow?
func AllowedParams() []string {
	return []string{"name", "url", "summary", "parent_id", "status"}
}

// Create a tag from database columns - used by query in creating objects
func NewWithColumns(cols map[string]interface{}) *Tag {

	tag := New()
	tag.Id = validate.Int(cols["id"])
	tag.CreatedAt = validate.Time(cols["created_at"])
	tag.UpdatedAt = validate.Time(cols["updated_at"])
	tag.ParentId = validate.Int(cols["parent_id"])
	tag.Status = validate.Int(cols["status"])
	tag.Name = validate.String(cols["name"])
	tag.Summary = validate.String(cols["summary"])
	tag.Url = validate.String(cols["url"])
	tag.DottedIds = validate.String(cols["dotted_ids"])

	return tag
}

// Set up a new tag with default values
func New() *Tag {
	tag := &Tag{}
	tag.Model.Init()
	tag.TableName = "tags"
	tag.Status = status.Draft
	tag.ParentId = 0
	tag.Url = ""
	tag.Name = ""
	tag.Summary = ""
	tag.DisplayOrder = 100000
	tag.DottedIds = ""
	return tag
}

// Insert a new tag
func Create(params map[string]string) (int64, error) {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Update/add some params by default
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// Create a new query relation referencing this model
func Query() *query.Query {
	return query.New("tags", "id")
}

// Create a new query for all models, setting a default order
func All() *query.Query {
	return Query().Order("updated_at desc, created_at desc, id desc")
}

func RootTags() *query.Query {
	return Query().Where("parent_id IS NULL OR parent_id = 0")
}

// A shortcut for the common where query on tags
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

func Ordered(q *query.Query) *query.Query {
	return q.Order("name asc")
}

// Request a single record by id in params
func Find(id int64) (*Tag, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Fetch all results for this query
func FindAll(q *query.Query) ([]*Tag, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results
	var tagList []*Tag
	for _, r := range results {
		tag := NewWithColumns(r)
		tagList = append(tagList, tag)
	}

	return tagList, nil
}

// Check these parameters conform to AcceptedParams, and pass validation
func validateParams(unsafeParams map[string]string) error {

	// Now check params are as we expect
	err := validate.Length(unsafeParams["id"], 0, -1)
	if err != nil {
		return err
	}

	return err
}

// Return the parent tag (if any)
func (m *Tag) Parent() *Tag {
	t, err := Find(m.ParentId)
	if err != nil {
		return nil
	}
	return t
}

// Update this tag
func (m *Tag) Update(params map[string]string) error {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return err
	}

	// Make sure updated_at is set to the current time
	params["updated_at"] = query.TimeString(time.Now().UTC())

	// Always regenerate dotted ids - we fetch all tags first to avoid db calls
	q := Query().Select("select id,parent_id from tags").Order("id asc")
	tagsList, err := FindAll(q)
	if err == nil {
		params["dotted_ids"] = m.CalculateDottedIds(tagsList)
	} else {
		return err
	}

	return Query().Where("id=?", m.Id).Update(params)
}

// Delete this tag
func (m *Tag) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// Return a list of tags suitable for parent options in a tag parent select
func (m *Tag) ParentTagOptions() []helpers.Option {

	options := make([]helpers.Option, 0)

	options = append(options, helpers.Option{0, "None"})

	q := Query().Order("name asc")
	tagsList, err := FindAll(q)
	if err == nil {
		for _, t := range tagsList {
			options = append(options, helpers.Option{t.Id, t.Name})
		}
	}

	return options
}

// Return a list of child tags by querying the database
func (m *Tag) Children() []*Tag {

	q := Query().Where("parent_id=?", m.Id).Order("name asc")

	// Fetch the tags
	tagsList, err := FindAll(q)
	if err != nil {
		fmt.Printf("Error fetching tag children %s", m.Name)
	}

	return tagsList
}

// Level returns our depth in the tag hierarchy as an int from 0 at root up
func (m *Tag) Level() int {
	if len(m.DottedIds) > 0 {
		return strings.Count(m.DottedIds, ".")
	}

	return 0

}

// CalculateDottedIds recalculates the dotted ids for this tag from parents (requires an array of all tag ids)
func (m *Tag) CalculateDottedIds(tags []*Tag) string {
	dottedIds := ""

	if m.ParentId != 0 {
		for _, tag := range tags {
			if tag.Id == m.ParentId {
				dottedIds = fmt.Sprintf("%s.%d", tag.CalculateDottedIds(tags), m.Id)
				break
			}
		}
	} else {
		dottedIds = fmt.Sprintf("%d", m.Id)
	}

	return dottedIds
}
