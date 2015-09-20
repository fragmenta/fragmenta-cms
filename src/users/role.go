package users

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"
)

// Role constants for this package
const (
	RoleAnon   = 0
	RoleReader = 20
	RoleEditor = 50
	RoleAdmin  = 100
)

// Anon returns true if this user is anon
func (m *User) Anon() bool {
	return m.Role == RoleAnon
}

// Patient returns true if this user is Patient
func (m *User) Reader() bool {
	return m.Role == RoleReader
}

// Expert returns true if this user is Expert
func (m *User) Editor() bool {
	return m.Role == RoleEditor
}

// Admin returns true if this user is Admin
func (m *User) Admin() bool {
	return m.Role == RoleAdmin
}

// RoleOptions returns an array of Role values for this model (embedders may override this and roledisplay to extend)
func (m *User) RoleOptions() []helpers.Option {
	var options []helpers.Option

	options = append(options, helpers.Option{Id: RoleReader, Name: "Reader"})
	options = append(options, helpers.Option{Id: RoleEditor, Name: "Editor"})
	options = append(options, helpers.Option{Id: RoleAdmin, Name: "Administrator"})

	return options
}

// RoleDisplay returns the string representation of the Role status
func (m *User) RoleDisplay() string {
	for _, o := range m.RoleOptions() {
		if o.Id == m.Role {
			return o.Name
		}
	}
	return ""
}

// Admins returns a query which finds all admin users
func Admins() *query.Query {
	return Query().Where("role=?", RoleAdmin).Order("name asc")
}

// Admins returns a query which finds all editor users
func Editors() *query.Query {
	return Query().Where("role=?", RoleEditor).Order("name asc")
}

// Readers returns a query  which finds all reader users
func Readers() *query.Query {
	return Query().Where("role=?", RoleReader).Order("name asc")
}
