package status

import (
	"testing"
)

// Resource embeds ResourceStatus
type resource struct {
	ResourceStatus
}

// TestOptions tests our options are functional when embedded in a resource.
func TestOptions(t *testing.T) {

	r := &resource{}

	options := r.StatusOptions()
	if len(options) < 0 {
		t.Fatalf("status: failed to get status options")
	}

	r.Status = Published
	if r.StatusDisplay() != "Published" {
		t.Fatalf("status: failed to get status published")
	}

}
