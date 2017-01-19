package mail

import (
	"errors"
	"fmt"

	"github.com/fragmenta/view"
)

// TODO - add more mail services, at present only sendgrid is supported
// Usage:
// email := mail.New(recipient)
// email.Subject = "blah"
// email.Body = blah
// mail.Send(email,context)

// Sender is the interface for our adapters for mail services.
type Sender interface {
	Send(email *Email) error
}

// Context defines a simple list of string:value pairs for mail templates.
type Context map[string]interface{}

// Production should be set to true in production environment.
var Production = false

// Service is the mail adapter to send with and should be set on startup.
var Service Sender

// Send the email using our default adapter and optional context.
func Send(email *Email, context Context) error {
	// If we have a template, render the email in that template
	if email.Body == "" && email.Template != "" {
		var err error
		email.Body, err = RenderTemplate(email, context)
		if err != nil {
			return err
		}
	}

	// If dev just log and return, don't send messages
	if !Production {
		fmt.Printf("#debug mail sent:%s\n", email)
		return nil
	}

	return Service.Send(email)
}

// RenderTemplate renders the email into its template with context.
func RenderTemplate(email *Email, context Context) (string, error) {
	if email.Template == "" || context == nil {
		return "", errors.New("mail: missing template or context")
	}

	view := view.NewWithPath("", nil)
	view.Layout(email.Layout)
	view.Template(email.Template)
	view.Context(context)
	body, err := view.RenderToStringWithLayout()
	if err != nil {
		return "", err
	}

	return body, nil
}
