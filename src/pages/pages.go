// Package pages represents the page resource on disk
package pages

import (
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Page handles saving and retreiving pages from the database
type Page struct {
	model.Model
	status.ModelStatus
	Url      string
	Name     string
	Summary  string
	Text     string
	Keywords string
	Template string
}

// AllowedParams returns an array of allowed param keys
func AllowedParams() []string {
	return []string{"url", "name", "summary", "text", "status", "template", "keywords"}
}

// NewWithColumns creates a new page instance and fills it with data from the database cols provided
func NewWithColumns(cols map[string]interface{}) *Page {

	page := New()
	page.Id = validate.Int(cols["id"])
	page.CreatedAt = validate.Time(cols["created_at"])
	page.UpdatedAt = validate.Time(cols["updated_at"])
	page.Status = validate.Int(cols["status"])
	page.Url = validate.String(cols["url"])
	page.Name = validate.String(cols["name"])
	page.Summary = validate.String(cols["summary"])
	page.Text = validate.String(cols["text"])
	page.Template = validate.String(cols["template"])
	page.Keywords = validate.String(cols["keywords"])

	return page
}

// New creates and initialises a new page instance
func New() *Page {
	page := &Page{}
	page.Model.Init()
	page.Status = status.Draft
	page.TableName = "pages"
	return page
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
func Find(id int64) (*Page, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll returns all results for this query
func FindAll(q *query.Query) ([]*Page, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results
	var pages []*Page
	for _, cols := range results {
		p := NewWithColumns(cols)
		pages = append(pages, p)
	}

	return pages, nil
}

// Query returns a new query for pages
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

// Published returns a query for all pages with status >= published
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}

// Where returns a Where query for pages with the arguments supplied
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Update sets the record in the database from params
func (m *Page) Update(params map[string]string) error {

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
func (m *Page) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// TemplateOptions returns an array of statuses for a status select
func (m *Page) TemplateOptions() []helpers.SelectableOption {
	var options []helpers.SelectableOption

	// We could instead iterate the possible views from view.Templates, that might be better
	options = append(options, helpers.SelectableOption{Name: "Page", Value: "pages/views/show.html.got"})
	options = append(options, helpers.SelectableOption{Name: "About Page", Value: "pages/views/show/about.html.got"})

	return options
}
