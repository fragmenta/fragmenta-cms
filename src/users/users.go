// Package users represents a logged in user
package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// User represents a user of the service
type User struct {
	model.Model
	status.ModelStatus

	Role              int64
	EncryptedPassword string
	Email             string
	Name              string
	Title             string
	Summary           string
	Text              string
	ImageID           int64
}

// AllowedParams returns an array of acceptable params in create/update
func AllowedParams() []string {
	return []string{"name", "summary", "email", "status", "role", "password", "text", "title", "image_id"}
}

// New create a user from database columns - used by query in creating objects
func NewWithColumns(cols map[string]interface{}) *User {

	user := New()
	user.Id = validate.Int(cols["id"])
	user.CreatedAt = validate.Time(cols["created_at"])
	user.UpdatedAt = validate.Time(cols["updated_at"])
	user.Status = validate.Int(cols["status"])
	user.Role = validate.Int(cols["role"])
	user.Name = validate.String(cols["name"])
	user.Summary = validate.String(cols["summary"])
	user.Email = validate.String(cols["email"])
	user.EncryptedPassword = validate.String(cols["encrypted_password"])
	user.Text = validate.String(cols["text"])
	user.Title = validate.String(cols["title"])
	user.ImageID = validate.Int(cols["image_id"])
	return user
}

// New sets up a new user with default values
func New() *User {
	user := &User{}
	user.Model.Init()
	user.TableName = "users"
	user.Status = status.Published
	user.Role = 0
	user.Name = ""
	user.Summary = ""
	user.Email = ""
	user.Text = "<h3>About</h3><p>About me</p>"
	return user
}

// Create inserts a new user
func Create(params map[string]string) (int64, error) {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Check that this user email is not already in use
	if len(params["email"]) > 0 {
		// Try to fetch a user by this email from the db - we don't allow duplicates
		count, err := Query().Where("email=?", params["email"]).Count()
		if err != nil {
			return 0, err
		}

		if count > 0 {
			return 0, errors.New("A username with this email already exists, sorry.")
		}

	}

	// Update/add some params by default
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// Query a new query relation referencing this model, optionally setting a default order
func Query() *query.Query {
	return query.New("users", "id")
}

// Where is a shortcut for the common where query on users
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Find requests a single record by id
func Find(id int64) (*User, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindEmail requests a single record by email
func FindEmail(email string) (*User, error) {
	result, err := Query().Where("email=?", email).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// First the first result for this query
func First(q *query.Query) (*User, error) {

	result, err := q.FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Fetch all results for this query
func FindAll(q *query.Query) ([]*User, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results
	var userList []*User
	for _, r := range results {
		user := NewWithColumns(r)
		userList = append(userList, user)
	}

	return userList, nil
}

// Exists checks whether a user email exists
func Exists(e string) bool {
	count, err := Query().Where("email=?", e).Count()
	if err != nil {
		return true // default to true on error
	}

	return (count > 0)
}

// validateParams these parameters conform to AcceptedParams, and pass validation
func validateParams(unsafeParams map[string]string) error {

	// Now check params are as we expect
	err := validate.Length(unsafeParams["id"], 0, -1)
	if err != nil {
		return err
	}

	err = validate.Length(unsafeParams["name"], 3, 100)
	if err != nil {
		return err
	}

	err = validate.Length(unsafeParams["email"], 3, 100)
	if err != nil {
		return err
	}

	// Password may be blank
	if len(unsafeParams["password"]) > 0 {
		// Report error for length between 0 and 5 chars
		err = validate.Length(unsafeParams["password"], 5, 100)
		if err != nil {
			return err
		}

		ep, err := auth.EncryptPassword(unsafeParams["password"])
		if err != nil {
			return err
		}
		unsafeParams["encrypted_password"] = ep

	}

	// Finally, always delete the password param
	delete(unsafeParams, "password")

	return err
}

// Update this user
func (m *User) Update(params map[string]string) error {

	// Remove params not in AllowedParams
	params = model.CleanParams(params, AllowedParams())

	err := validateParams(params)
	if err != nil {
		return err
	}

	// Make sure updated_at is set to the current time
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// Destroy this user
func (m *User) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// URLShow returns the url for this user
func (m *User) URLShow() string {
	return fmt.Sprintf("/users/%d-%s", m.Id, m.ToSlug(m.Name))
}

// SelectName returns the name which should be used for select options
func (m *User) SelectName() string {
	return m.Name
}

// Keywords returns meta keywords for this user
func (m *User) Keywords() string {
	return fmt.Sprintf("%s %s", m.Name, m.Summary)
}
