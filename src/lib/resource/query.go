package resource

import (
	"time"

	"github.com/fragmenta/query"
)

// Query creates a new query relation referencing this specific resource by id.
func (r *Base) Query() *query.Query {
	return query.New(r.Table(), r.PrimaryKey()).Where("id=?", r.ID)
}

// ValidateParams allows only those params by AllowedParams()
// to perform more sophisticated validation override it.
func (r *Base) ValidateParams(params map[string]string, allowed []string) map[string]string {

	for k := range params {
		paramAllowed := false
		for _, v := range allowed {
			if k == v {
				paramAllowed = true
			}
		}
		if !paramAllowed {
			delete(params, k)
		}
	}
	return params
}

// Create inserts a new database record and returns the id or an error
func (r *Base) Create(params map[string]string) (int64, error) {

	// Make sure updated_at and created_at are set to the current time
	now := query.TimeString(time.Now().UTC())
	params["created_at"] = now
	params["updated_at"] = now

	// Insert a record into the database
	id, err := query.New(r.Table(), r.PrimaryKey()).Insert(params)
	return id, err
}

// Update the database record for this resource with the given params.
func (r *Base) Update(params map[string]string) error {

	// Make sure updated_at is set to the current time
	now := query.TimeString(time.Now().UTC())
	params["updated_at"] = now

	return r.Query().Update(params)
}

// Destroy deletes this resource by removing the database record.
func (r *Base) Destroy() error {
	return r.Query().Delete()
}
