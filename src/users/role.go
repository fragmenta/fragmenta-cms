package users

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
)

// This file contains functions related to authorisation and roles.

// User roles
const (
	Anon   = 0
	Editor = 10
	Reader = 20
	Admin  = 100
)

// RoleOptions returns an array of Role values for this model (embedders may override this and roledisplay to extend)
func (u *User) RoleOptions() []helpers.Option {
	var options []helpers.Option

	options = append(options, helpers.Option{Id: Reader, Name: "Reader"})
	options = append(options, helpers.Option{Id: Editor, Name: "Editor"})
	options = append(options, helpers.Option{Id: Admin, Name: "Administrator"})

	return options
}

// RoleDisplay returns the string representation of the Role status
func (u *User) RoleDisplay() string {
	for _, o := range u.RoleOptions() {
		if o.Id == u.Role {
			return o.Name
		}
	}
	return ""
}

// Anon returns true if this user is not a logged in user.
func (u *User) Anon() bool {
	return u.Role == Anon || u.ID == 0
}

// Admin returns true if this user is an Admin.
func (u *User) Admin() bool {
	return u.Role == Admin
}

// Reader returns true if this user is an Reader.
func (u *User) Reader() bool {
	return u.Role == Reader
}

// Admins returns a query which finds all admin users
func Admins() *query.Query {
	return Query().Where("role=?", Admin).Order("name asc")
}

// Editors returns a query which finds all editor users
func Editors() *query.Query {
	return Query().Where("role=?", Editor).Order("name asc")
}

// Readers returns a query  which finds all reader users
func Readers() *query.Query {
	return Query().Where("role=?", Reader).Order("name asc")
}

// can.User interface

// RoleID returns the user role for auth.
func (u *User) RoleID() int64 {
	if u == nil {
		return Anon
	}
	return u.Role
}

// UserID returns the user id for auth.
func (u *User) UserID() int64 {
	if u == nil {
		return 0
	}
	return u.ID
}

// MockAnon returns a mock user for testing with Role Anon.
func MockAnon() *User {
	return &User{Role: Anon, Email: "anon@example.com"}
}

// MockAdmin returns a mock user for testing with Role Admin.
func MockAdmin() *User {
	return &User{Role: Admin, Email: "admin@example.com", Base: resource.Base{ID: 1}}
}
