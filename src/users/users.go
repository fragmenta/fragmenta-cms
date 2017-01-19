// Package users represents the user resource
package users

import (
	"time"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// User handles saving and retreiving users from the database
type User struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	// Authorisation
	Role int64

	// Authentication
	PasswordHash       string
	PasswordResetToken string
	PasswordResetAt    time.Time

	// User details
	Email   string
	Name    string
	Title   string
	Summary string
	Text    string
	ImageID int64
}
