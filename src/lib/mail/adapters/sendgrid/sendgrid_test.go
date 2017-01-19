package sendgrid

import (
	"testing"

	"github.com/fragmenta/fragmenta-cms/src/lib/mail"
)

// TestSendGrid tests we can create a new sendgrid connection
// without mocking we can't really test an actual send
func TestSendGrid(t *testing.T) {

	// TODO add some proper tests here.

	s := New("from@example.com", "secret")
	if s.secret != "secret" {
		t.Errorf("sendgrid: failed to set up")
	}

	email := mail.New("example@example.com")
	err := s.Send(email)
	if err == nil {
		t.Errorf("sendgrid: failed to error on bad emails")
	}

}
