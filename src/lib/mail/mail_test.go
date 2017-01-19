package mail

import (
	"testing"

	"github.com/fragmenta/view"
)

// TestMail tests that mail formats properly in dev mode
func TestMail(t *testing.T) {

	// In order to test, we rely on the view pkg being set up
	err := view.LoadTemplatesAtPaths([]string{"../.."}, view.Helpers)
	if err != nil {
		t.Errorf("mail: failed to load views")
	}

	context := Context{
		"msg": "hello world",
	}

	recipient := "recipient@example.com"
	email := New(recipient)
	email.ReplyTo = "sender@example.com"
	email.Subject = "sub"
	email.Body = "<h1>Body</h1>"
	Send(email, context)

	// Try render
	email.Body, err = RenderTemplate(email, context)
	if err != nil {
		t.Errorf("mail: failed to render message :%s", err)
	}

}
