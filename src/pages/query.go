package pages

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "pages"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"status", "author_id", "keywords", "name", "status", "summary", "template", "text", "url"}
}

// NewWithColumns creates a new page instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Pages {

	page := New()
	page.ID = resource.ValidateInt(cols["id"])
	page.CreatedAt = resource.ValidateTime(cols["created_at"])
	page.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	page.Status = resource.ValidateInt(cols["status"])
	page.AuthorID = resource.ValidateInt(cols["author_id"])
	page.Keywords = resource.ValidateString(cols["keywords"])
	page.Name = resource.ValidateString(cols["name"])
	page.Status = resource.ValidateInt(cols["status"])
	page.Summary = resource.ValidateString(cols["summary"])
	page.Template = resource.ValidateString(cols["template"])
	page.Text = resource.ValidateString(cols["text"])
	page.URL = resource.ValidateString(cols["url"])

	return page
}

// New creates and initialises a new page instance.
func New() *Pages {
	page := &Pages{}
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()
	page.TableName = TableName
	page.KeyName = KeyName
	page.Status = status.Draft
	page.Template = "pages/views/templates/default.html.got"
	return page
}

// FindFirst fetches a single page record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Pages, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single page record from the database by id.
func Find(id int64) (*Pages, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all page records matching this query from the database.
func FindAll(q *query.Query) ([]*Pages, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results
	var pages []*Pages
	for _, cols := range results {
		p := NewWithColumns(cols)
		pages = append(pages, p)
	}

	return pages, nil
}

// Query returns a new query for pages with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for pages with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all pages with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
