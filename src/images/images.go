// Package images represents the image resource
package images

import (
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Image handles saving and retreiving images from the database
type Image struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	AuthorID int64
	Name     string
	Path     string
	Sort     int64
	Status   int64
}
