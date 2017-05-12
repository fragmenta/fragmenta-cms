package images

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "images"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"status", "author_id", "name", "path", "sort", "status"}
}

// NewWithColumns creates a new image instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Image {

	image := New()
	image.ID = resource.ValidateInt(cols["id"])
	image.CreatedAt = resource.ValidateTime(cols["created_at"])
	image.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	image.Status = resource.ValidateInt(cols["status"])
	image.AuthorID = resource.ValidateInt(cols["author_id"])
	image.Name = resource.ValidateString(cols["name"])
	image.Path = resource.ValidateString(cols["path"])
	image.Sort = resource.ValidateInt(cols["sort"])
	image.Status = resource.ValidateInt(cols["status"])

	return image
}

// New creates and initialises a new image instance.
func New() *Image {
	image := &Image{}
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()
	image.TableName = TableName
	image.KeyName = KeyName
	image.Status = status.Draft
	return image
}

// FindFirst fetches a single image record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Image, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single image record from the database by id.
func Find(id int64) (*Image, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all image records matching this query from the database.
func FindAll(q *query.Query) ([]*Image, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of images constructed from the results
	var images []*Image
	for _, cols := range results {
		p := NewWithColumns(cols)
		images = append(images, p)
	}

	return images, nil
}

// Query returns a new query for images with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for images with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all images with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
