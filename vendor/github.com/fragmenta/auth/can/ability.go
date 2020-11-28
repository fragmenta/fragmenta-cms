// Package can implements basic role-based permissions for golang
// - controlling who can.Do certain actions for a given database table.
package can

import (
	"errors"
	"fmt"
	"sync"
)

// abilities is an array of abilities
var abilities []*Ability

// mu protects the list of abilities during access
var mu sync.RWMutex

// Authorise adds this ability to the list of abilities for this role.
// Usage: can.Authorise(role.Admin, can.ManageResource, "pages")
func Authorise(role int64, v Verb, id string) {
	ability := &Ability{role: role, verb: v, identifier: id, ownership: false}
	add(ability)
}

// AuthoriseOwner adds this ability to the list of abilities for this role
// for resources owned by this user.
// Usage: can.AuthoriseOwner(role.Reader, can.ShowResource, "pages")
func AuthoriseOwner(role int64, v Verb, id string) {
	ability := &Ability{role: role, verb: v, identifier: id, ownership: true}
	add(ability)
}

// add adds this ability
func add(a *Ability) {
	mu.Lock()
	abilities = append(abilities, a)
	mu.Unlock()
}

// Ability represents an authorisation for an action for a given role
type Ability struct {
	ownership  bool
	role       int64
	verb       Verb
	identifier string
}

// Allow returns an error if the action is not allowed, or nil if it is
func (a *Ability) Allow(v Verb, r Resource, u User) error {

	// Fail if user role doesn't match
	if u == nil || a.role != u.RoleID() {
		return errors.New("can: role not authorised")
	}

	// Fail if resource id doesn't match
	if a.identifier != Anything && a.identifier != r.ResourceID() {
		return errors.New("can: resource not authorised")
	}

	// Check for verb match, fail if no match
	if a.verb != ManageResource && a.verb != v {
		return errors.New("can: action not authorised")
	}

	// If we have an ability which doesn't require ownership, return now
	if !a.CheckOwner() {
		return nil
	}

	// Now check ownership
	if r == nil || !r.OwnedBy(u.UserID()) {
		return errors.New("can: action not authorised")
	}

	return nil
}

// CheckOwner returns true if this ability should check ownership
func (a *Ability) CheckOwner() bool {
	// If the verb is to create or list, we can do no ownership check
	if a.verb == CreateResource || a.verb == ListResource {
		return false
	}
	// If the resource is anything, we do not check ownership
	if a.identifier == Anything {
		return false
	}
	// If the ability does not require ownership, return false
	return a.ownership
}

// String returns a string description of this ability.
func (a *Ability) String() string {
	return fmt.Sprintf("%v %d can %v on %s\n", a.ownership, a.role, a.verb, a.identifier)
}
