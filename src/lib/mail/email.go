package mail

import (
	"fmt"
)

// Email represents an email to be sent.
type Email struct {
	Recipients []string
	ReplyTo    string
	Subject    string
	Body       string
	Template   string
	Layout     string
}

// New returns a new email with the default tenplates and the given recipient.
func New(r string) *Email {
	e := &Email{
		Layout:   "lib/mail/views/layout.html.got",
		Template: "lib/mail/views/template.html.got",
	}
	e.Recipients = append(e.Recipients, r)
	return e
}

// String returns a formatted string representation for debug.
func (e *Email) String() string {
	return fmt.Sprintf("email to:%v from:%s subject:%s\n\n%s", e.Recipients, e.ReplyTo, e.Subject, e.Body)
}

// Invalid returns true if this email is not ready to send.
func (e *Email) Invalid() bool {
	return (e.ReplyTo == "" || e.Subject == "" || e.Body == "")
}
