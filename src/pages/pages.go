// Package pages represents the page resource
package pages

import (
	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Page handles saving and retreiving pages from the database
type Page struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	AuthorID int64
	Keywords string
	Name     string
	Status   int64
	Summary  string
	Template string
	Text     string
	URL      string
}

// ShowTemplate returns the default template if none is set, or the template selected
func (p *Page) ShowTemplate() string {
	if p.Template == "" {
		return "pages/views/templates/default.html.got"
	}
	return p.Template
}
