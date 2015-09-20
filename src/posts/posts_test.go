// Tests for the posts package
package posts

import (
	"testing"
)

// Log a failure message, given msg, expected and result
func fail(t *testing.T, msg string, expected interface{}, result interface{}) {
	t.Fatalf("\n------FAILURE------\nTest failed: %s expected:%v result:%v", msg, expected, result)
}

// Test create of Post
func TestCreatePost(t *testing.T) {

}

// Test update of Post
func TestUpdatePost(t *testing.T) {

}
