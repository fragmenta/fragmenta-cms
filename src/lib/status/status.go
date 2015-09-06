package status

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/view/helpers"
)

// ModelStatus adds status to a model - this may in future be removed and moved into apps as it is frequently modified
type ModelStatus struct {
	Status int64
}

// Status values
// If these need to be substantially modified for a particular model,
// it may be better to move this into the model package concerned and modify as required
const (
	Draft       = 0
	Final       = 10
	Suspended   = 11
	Unavailable = 12
	Published   = 100
	Featured    = 101
)

// StatusOptions returns an array of statuses for a status select
func (m *ModelStatus) StatusOptions() []helpers.Option {
	var options []helpers.Option

	options = append(options, helpers.Option{Id: Draft, Name: "Draft"})
	options = append(options, helpers.Option{Id: Final, Name: "Final"})
	options = append(options, helpers.Option{Id: Suspended, Name: "Suspended"})
	options = append(options, helpers.Option{Id: Published, Name: "Published"})

	return options
}

// StatusDisplay returns a string representation of the model status
func (m *ModelStatus) StatusDisplay() string {
	for _, o := range m.StatusOptions() {
		if o.Id == m.Status {
			return o.Name
		}
	}
	return ""
}

// Model status

// IsDraft returns true if the status is Draft
func (m *ModelStatus) IsDraft() bool {
	return m.Status == Draft
}

// IsFinal returns true if the status is Final
func (m *ModelStatus) IsFinal() bool {
	return m.Status == Final
}

// IsSuspended returns true if the status is Suspended
func (m *ModelStatus) IsSuspended() bool {
	return m.Status == Suspended
}

// IsUnavailable returns true if the status is unavailable
func (m *ModelStatus) IsUnavailable() bool {
	return m.Status == Unavailable
}

// IsPublished returns true if the status is published *or over*
func (m *ModelStatus) IsPublished() bool {
	return m.Status >= Published // NB >=
}

// IsFeatured returns true if the status is featured
func (m *ModelStatus) IsFeatured() bool {
	return m.Status == Featured
}

// CHAINABLE FINDER FUNCTIONS
// Apply with query.Apply(status.WherePublished) etc
// Or define on your own models instead...

// WhereDraft modifies the given query to select status draft
func WhereDraft(q *query.Query) *query.Query {
	return q.Where("status = ?", Draft)
}

// WhereFinal modifies the given query to select status Final
func WhereFinal(q *query.Query) *query.Query {
	return q.Where("status = ?", Final)
}

// WhereSuspended modifies the given query to select status Suspended
func WhereSuspended(q *query.Query) *query.Query {
	return q.Where("status = ?", Suspended)
}

// WhereFeatured modifies the given query to select status Featured
func WhereFeatured(q *query.Query) *query.Query {
	return q.Where("status = ?", Featured)
}

// WherePublished modifies the given query to select status Published
func WherePublished(q *query.Query) *query.Query {
	return q.Where("status >= ?", Published)
}

// Null modifies the given query to select records with null status
func Null(q *query.Query) *query.Query {
	return q.Where("status IS NULL")
}

// NotNull modifies the given query to select records which do not have null status
func NotNull(q *query.Query) *query.Query {
	return q.Where("status IS NOT NULL")
}

// Order modifies the given query to order records by status
func Order(q *query.Query) *query.Query {
	return q.Order("status desc")
}
