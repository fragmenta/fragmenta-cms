// Tests for the users package
package users

import (
	"testing"
)

// Log a failure message, given msg, expected and result
func fail(t *testing.T, msg string, expected interface{}, result interface{}) {
	t.Fatalf("Test failed: %s expected:%v result:%v", msg, expected, result)
}

// Test
func TestUser(t *testing.T) {

}
