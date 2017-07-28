// Package redirects represents the redirect resource
package redirects

import (
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Redirect handles saving and retreiving redirects from the database
type Redirect struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	NewURL string
	OldURL string
}
