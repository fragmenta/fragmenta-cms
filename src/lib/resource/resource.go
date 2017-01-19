// Package resource provides some shared behaviour for resources, and basic CRUD and URL helpers.
package resource

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/fragmenta/auth"
)

// Base defines shared fields and behaviour for resources.
type Base struct {
	// ID is the default primary key of the resource.
	ID int64

	// CreatedAt stores the creation time of the resource.
	CreatedAt time.Time

	// UpdatedAt stores the last update time of the resource.
	UpdatedAt time.Time

	// TableName is used for database queries and urls.
	TableName string

	// KeyName is used for database queries as the primary key.
	KeyName string
}

// String returns a string representation of the resource
func (r *Base) String() string {
	return fmt.Sprintf("%s/%d", r.TableName, r.ID)
}

// Queryable interface

// Table returns the table name for this object
func (r *Base) Table() string {
	return r.TableName
}

// PrimaryKey returns the id for primary key by default - used by query
func (r *Base) PrimaryKey() string {
	return r.KeyName
}

// PrimaryKeyValue returns the unique id
func (r *Base) PrimaryKeyValue() int64 {
	return r.ID
}

// Selectable interface

// SelectName returns our name for select menus
func (r *Base) SelectName() string {
	return fmt.Sprintf("%s-%d", r.TableName, r.ID)
}

// SelectValue returns our value for select options
func (r *Base) SelectValue() string {
	return fmt.Sprintf("%d", r.ID)
}

// Cacheable interface

// CacheKey generates a cache key for this resource
// based on the TableName, ID and UpdatedAt
func (r *Base) CacheKey() string {
	key := []byte(fmt.Sprintf("%s/%d/%s", r.TableName, r.ID, r.UpdatedAt))
	hash := sha256.Sum256(key)
	return auth.BytesToHex(hash[:32])
}

// can.Resource interface

// OwnedBy returns true if the user id passed in owns this resource.
func (r *Base) OwnedBy(uid int64) bool {
	return false
}

// ResourceID returns a key unique to this resource (we use table).
func (r *Base) ResourceID() string {
	return r.TableName
}
