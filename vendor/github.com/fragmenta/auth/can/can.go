package can

import (
	"fmt"
)

// User defines the interface for users which must have numeric roles
type User interface {
	RoleID() int64 // for role check
	UserID() int64 // for ownership check
}

// Resource defines the interface for resources
type Resource interface {
	OwnedBy(int64) bool // for ownership check, passed a UserID
	ResourceID() string // for check against abilities registered on this resource
}

// Verb represents the action taken on resources
type Verb int

// Verbs used to authorise actions on resources.
// Manages allows any action on a resource,
// and all verbs after Creates check ownership of the resource with OwnedBy().
const (
	ManageResource = iota
	ListResource   // Does not check ownership
	CreateResource // Does not check ownership
	ShowResource
	UpdateResource
	DestroyResource
)

// Resource identifier used to short-circuit checks on resource identity in conjuction with ManageResource
const (
	Anything = "*" // Allow actions on any resource
)

// Do returns an error if this action is not allowed, or nil if it is allowed
func Do(v Verb, r Resource, u User) error {

	// Check abilities for a match
	mu.RLock()
	for _, a := range abilities {

		// If no err, return nil to signify success
		if a.Allow(v, r, u) == nil {
			return nil
		}
	}
	mu.RUnlock()

	// If we reach here, no matching authorisation was found - note u may be nil
	return fmt.Errorf("can: no authorisation for action:%v %v %v", v, r, u)
}

// The following are wrapper functions for can.Do to provide a more elegant interface
// i.e. calling can.Manage(u,r)

// Manage returns an error if all actions are not authorised for this user
func Manage(r Resource, u User) error {
	return Do(ManageResource, r, u)
}

// Create returns an error if this action is not authorised for this user
func Create(r Resource, u User) error {
	return Do(CreateResource, r, u)
}

// List returns an error if this action is not authorised for this user
func List(r Resource, u User) error {
	return Do(ListResource, r, u)
}

// Show returns an error if this action is not authorised for this user
func Show(r Resource, u User) error {
	return Do(ShowResource, r, u)
}

// Update returns an error if this action is not authorised for this user
func Update(r Resource, u User) error {
	return Do(UpdateResource, r, u)
}

// Destroy returns an error if this action is not authorised for this user
func Destroy(r Resource, u User) error {
	return Do(DestroyResource, r, u)
}
