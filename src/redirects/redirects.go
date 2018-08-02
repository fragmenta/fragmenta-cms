// Package redirects represents the redirect resource
package redirects

import (
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

// Redirect handles saving and retreiving redirects from the database
type Redirect struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	NewURL string
	OldURL string
}
