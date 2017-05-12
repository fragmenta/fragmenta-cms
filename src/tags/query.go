package tags

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "tags"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"status", "dotted_ids", "name", "parent_id", "sort", "status", "summary", "url"}
}

// NewWithColumns creates a new tag instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Tag {

	tag := New()
	tag.ID = resource.ValidateInt(cols["id"])
	tag.CreatedAt = resource.ValidateTime(cols["created_at"])
	tag.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	tag.Status = resource.ValidateInt(cols["status"])
	tag.DottedIDs = resource.ValidateString(cols["dotted_ids"])
	tag.Name = resource.ValidateString(cols["name"])
	tag.ParentID = resource.ValidateInt(cols["parent_id"])
	tag.Sort = resource.ValidateInt(cols["sort"])
	tag.Status = resource.ValidateInt(cols["status"])
	tag.Summary = resource.ValidateString(cols["summary"])
	tag.URL = resource.ValidateString(cols["url"])

	return tag
}

// New creates and initialises a new tag instance.
func New() *Tag {
	tag := &Tag{}
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	tag.TableName = TableName
	tag.KeyName = KeyName
	tag.Status = status.Draft
	return tag
}

// FindFirst fetches a single tag record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Tag, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single tag record from the database by id.
func Find(id int64) (*Tag, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all tag records matching this query from the database.
func FindAll(q *query.Query) ([]*Tag, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of tags constructed from the results
	var tags []*Tag
	for _, cols := range results {
		p := NewWithColumns(cols)
		tags = append(tags, p)
	}

	return tags, nil
}

// Query returns a new query for tags with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for tags with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all tags with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
