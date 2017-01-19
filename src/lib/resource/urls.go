package resource

import (
	"fmt"
)

// IndexURL returns the index url for this model - /table
func (r *Base) IndexURL() string {
	return fmt.Sprintf("/%s", r.TableName)
}

// CreateURL returns the create url for this model /table/create
func (r *Base) CreateURL() string {
	return fmt.Sprintf("/%s/create", r.TableName)
}

// UpdateURL returns the update url for this model /table/id/update
func (r *Base) UpdateURL() string {
	return fmt.Sprintf("/%s/%d/update", r.TableName, r.ID)
}

// DestroyURL returns the destroy url for this model /table/id/destroy
func (r *Base) DestroyURL() string {
	return fmt.Sprintf("/%s/%d/destroy", r.TableName, r.ID)
}

// ShowURL returns the show url for this model /table/id
func (r *Base) ShowURL() string {
	return fmt.Sprintf("/%s/%d", r.TableName, r.ID)
}

// PublicURL returns the canonical url for showing this resource
// usually this will differ in using the name as a slug
func (r *Base) PublicURL() string {
	return fmt.Sprintf("/%s/%d", r.TableName, r.ID)
}
