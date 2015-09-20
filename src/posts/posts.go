// Package posts represents the post resource
package posts

import (
	"fmt"
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Post handles saving and retreiving posts from the database
type Post struct {
	model.Model
	status.ModelStatus
	AuthorId int64
	Name     string
	Summary  string
	Text     string
}

// AllowedParams returns an array of allowed param keys
func AllowedParams() []string {
	return []string{"status", "author_id", "name", "summary", "text"}
}

// NewWithColumns creates a new post instance and fills it with data from the database cols provided
func NewWithColumns(cols map[string]interface{}) *Post {

	post := New()
	post.Id = validate.Int(cols["id"])
	post.CreatedAt = validate.Time(cols["created_at"])
	post.UpdatedAt = validate.Time(cols["updated_at"])
	post.Status = validate.Int(cols["status"])
	post.AuthorId = validate.Int(cols["author_id"])
	post.Name = validate.String(cols["name"])
	post.Summary = validate.String(cols["summary"])
	post.Text = validate.String(cols["text"])

	return post
}

// New creates and initialises a new post instance
func New() *Post {
	post := &Post{}
	post.Model.Init()
	post.Status = status.Draft
	post.TableName = "posts"
	post.Text = "<section class=\"padded\"><h1>Title</h1><p>Text</p></section>"
	return post
}

// Create inserts a new record in the database using params, and returns the newly created id
func Create(params map[string]string) (int64, error) {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	// Check params for invalid values
	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Update date params
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// validateParams checks these params pass validation checks
func validateParams(params map[string]string) error {

	// Now check params are as we expect
	err := validate.Length(params["id"], 0, -1)
	if err != nil {
		return err
	}
	err = validate.Length(params["name"], 0, 255)
	if err != nil {
		return err
	}

	return err
}

// Find returns a single record by id in params
func Find(id int64) (*Post, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll returns all results for this query
func FindAll(q *query.Query) ([]*Post, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of posts constructed from the results
	var posts []*Post
	for _, cols := range results {
		p := NewWithColumns(cols)
		posts = append(posts, p)
	}

	return posts, nil
}

// Query returns a new query for posts
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

// Published returns a query for all posts with status >= published
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}

// Where returns a Where query for posts with the arguments supplied
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Update sets the record in the database from params
func (m *Post) Update(params map[string]string) error {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	// Check params for invalid values
	err := validateParams(params)
	if err != nil {
		return err
	}

	// Update date params
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// Destroy removes the record from the database
func (m *Post) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// URLShow returns an url with a slug
func (m *Post) URLShow() string {
	return fmt.Sprintf("/posts/%d-%s", m.Id, m.ToSlug(m.Name))
}
