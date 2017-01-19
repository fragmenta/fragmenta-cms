package status

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"
)

// Status values valid in the status field added with status.ResourceStatus.
const (
	None      = 0
	Draft     = 1
	Suspended = 50
	Published = 100
)

// ResourceStatus adds a status field to resources.
type ResourceStatus struct {
	Status int64
}

// WherePublished modifies the given query to select status greater than published.
// Note this selects >= Published.
func WherePublished(q *query.Query) *query.Query {
	return q.Where("status >= ?", Published)
}

// Options returns an array of statuses for a status select.
func Options() []helpers.Option {
	var options []helpers.Option

	options = append(options, helpers.Option{Id: Draft, Name: "Draft"})
	options = append(options, helpers.Option{Id: Suspended, Name: "Suspended"})
	options = append(options, helpers.Option{Id: Published, Name: "Published"})

	return options
}

// OptionsAll returns a list of options starting with a None option using the name passed in,
// which is useful for filter menus filtering on status.
func OptionsAll(name string) []helpers.Option {
	options := Options()
	return append(options, helpers.Option{Id: None, Name: name})
}

// StatusOptions returns an array of statuses for a status select for this resource.
func (r *ResourceStatus) StatusOptions() []helpers.Option {
	return Options()
}

// StatusDisplay returns a string representation of the model status.
func (r *ResourceStatus) StatusDisplay() string {
	for _, o := range r.StatusOptions() {
		if o.Id == r.Status {
			return o.Name
		}
	}
	return ""
}
