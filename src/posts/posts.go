// Package posts represents the post resource
package posts

import (
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
	Status   int64
	Summary  string
	Template string
	Text     string
}

func (p *Post) StatusDisplay() string {
	for _, o := range p.StatusOptions() {
		if o.Id == p.Status {
			return o.Name
		}
	}
	return ""
}
