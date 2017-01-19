package users

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "users"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"name", "summary", "email", "status", "role", "password", "text", "title", "image_id"}
}

// NewWithColumns creates a new user instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *User {

	user := New()
	user.TableName = TableName
	user.KeyName = KeyName
	user.ID = resource.ValidateInt(cols["id"])
	user.CreatedAt = resource.ValidateTime(cols["created_at"])
	user.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	user.Status = resource.ValidateInt(cols["status"])
	user.Email = resource.ValidateString(cols["email"])
	user.Name = resource.ValidateString(cols["name"])
	user.Role = resource.ValidateInt(cols["role"])
	user.Summary = resource.ValidateString(cols["summary"])
	user.Text = resource.ValidateString(cols["text"])
	user.Title = resource.ValidateString(cols["title"])
	user.ImageID = resource.ValidateInt(cols["image_id"])
	user.PasswordHash = resource.ValidateString(cols["password_hash"])
	user.PasswordResetToken = resource.ValidateString(cols["password_reset_token"])
	user.PasswordResetAt = resource.ValidateTime(cols["password_reset_at"])

	return user
}

// New creates and initialises a new user instance.
func New() *User {
	user := &User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.TableName = TableName
	user.KeyName = KeyName
	user.Status = status.Draft
	return user
}

// FindFirst fetches a single user record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*User, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single user record from the database by id.
func Find(id int64) (*User, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all user records matching this query from the database.
func FindAll(q *query.Query) ([]*User, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of users constructed from the results
	var users []*User
	for _, cols := range results {
		p := NewWithColumns(cols)
		users = append(users, p)
	}

	return users, nil
}

// Query returns a new query for users with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for users with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all users with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
