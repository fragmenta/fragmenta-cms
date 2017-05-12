package posts

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "posts"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"status", "author_id", "keywords", "name", "status", "summary", "template", "text"}
}

// NewWithColumns creates a new post instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Posts {

	post := New()
	post.ID = resource.ValidateInt(cols["id"])
	post.CreatedAt = resource.ValidateTime(cols["created_at"])
	post.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	post.Status = resource.ValidateInt(cols["status"])
	post.AuthorID = resource.ValidateInt(cols["author_id"])
	post.Keywords = resource.ValidateString(cols["keywords"])
	post.Name = resource.ValidateString(cols["name"])
	post.Status = resource.ValidateInt(cols["status"])
	post.Summary = resource.ValidateString(cols["summary"])
	post.Template = resource.ValidateString(cols["template"])
	post.Text = resource.ValidateString(cols["text"])

	return post
}

// New creates and initialises a new post instance.
func New() *Posts {
	post := &Posts{}
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.TableName = TableName
	post.KeyName = KeyName
	post.Status = status.Draft
	post.Template = "posts/views/templates/default.html.got"
	return post
}

// FindFirst fetches a single post record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Posts, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single post record from the database by id.
func Find(id int64) (*Posts, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all post records matching this query from the database.
func FindAll(q *query.Query) ([]*Posts, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of posts constructed from the results
	var posts []*Posts
	for _, cols := range results {
		p := NewWithColumns(cols)
		posts = append(posts, p)
	}

	return posts, nil
}

// Query returns a new query for posts with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for posts with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all posts with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
