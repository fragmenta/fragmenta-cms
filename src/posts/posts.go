// Package posts represents the post resource
package posts

import (
	"github.com/fragmenta/view/helpers"

	"github.com/fragmenta/fragmenta-cms/src/lib/resource"
	"github.com/fragmenta/fragmenta-cms/src/lib/status"
)

// Post handles saving and retreiving posts from the database
type Post struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	AuthorID int64
	Keywords string
	Name     string
	Summary  string
	Template string
	Text     string
}

// ShowTemplate returns the default template if none is set, or the template selected
func (p *Post) ShowTemplate() string {
	if p.Template == "" {
		return "posts/views/templates/default.html.got"
	}
	return p.Template
}

// TemplateOptions provides a set of options for the templates menu
// ids are indexes into the templates array above
func (p *Post) TemplateOptions() []helpers.Selectable {
	var options []helpers.Selectable

	options = append(options, helpers.SelectableOption{Value: "posts/views/templates/default.html.got", Name: "Default"})

	return options
}
