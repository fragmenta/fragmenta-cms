// Package tags represents the tag resource
package tags

import (
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Tags handles saving and retreiving tags from the database
type Tags struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	DottedIDs string
	Name      string
	ParentID  int64
	Sort      int64
	Status    int64
	Summary   string
	URL       string
}
