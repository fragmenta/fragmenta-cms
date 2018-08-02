package redirects

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

const (
	// TableName is the database table for this resource
	TableName = "redirects"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "updated_at desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"new_url", "old_url"}
}

// NewWithColumns creates a new redirect instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Redirect {

	redirect := New()
	redirect.ID = resource.ValidateInt(cols["id"])
	redirect.CreatedAt = resource.ValidateTime(cols["created_at"])
	redirect.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	redirect.NewURL = resource.ValidateString(cols["new_url"])
	redirect.OldURL = resource.ValidateString(cols["old_url"])

	return redirect
}

// New creates and initialises a new redirect instance.
func New() *Redirect {
	redirect := &Redirect{}
	redirect.CreatedAt = time.Now()
	redirect.UpdatedAt = time.Now()
	redirect.TableName = TableName
	redirect.KeyName = KeyName
	return redirect
}

// FindFirst fetches a single redirect record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Redirect, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single redirect record from the database by id.
func Find(id int64) (*Redirect, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all redirect records matching this query from the database.
func FindAll(q *query.Query) ([]*Redirect, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of redirects constructed from the results
	var redirects []*Redirect
	for _, cols := range results {
		p := NewWithColumns(cols)
		redirects = append(redirects, p)
	}

	return redirects, nil
}

// Query returns a new query for redirects with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for redirects with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}
