package mail

import (
	"github.com/fragmenta/view"
	"github.com/sendgrid/sendgrid-go"
)

// TODO - instead of package variables, use New() mailer and put the variables in the mailer?
// Try to be consistent about use of pkg variables

// The Mail service user (must be set before first sending)
var user string

// The Mail service secret key/password (must be set before first sending)
var secret string

var from string

// Setup sets the user and secret for use in sending mail (possibly later we should have a config etc)
func Setup(u string, s string, f string) {
	user = u
	secret = s
	from = f
}

// Send sends mail
func Send(recipients []string, subject string, template string, context map[string]interface{}) error {

	// At present we use sendgrid, we may later allow many mail services to be used
	//	sg := sendgrid.NewSendGridClient(user, secret)
	sg := sendgrid.NewSendGridClientWithApiKey(secret)

	message := sendgrid.NewMail()
	message.SetFrom(from)
	message.AddTos(recipients)
	message.SetSubject(subject)

	// Hack for now, consider options TODO
	if context["reply_to"] != nil {
		replyTo := context["reply_to"].(string)
		message.SetReplyTo(replyTo)
	}

	// Load the template, and substitute using context
	// We should possibly set layout from caller too?
	view := view.New(&renderContext{})
	view.Template(template)
	view.Context(context)

	// We have a panic: runtime error: invalid memory address or nil pointer dereference
	// because of reloading templates on the fly I think
	// github.com/fragmenta/view/parser/template.html.go:90
	html, err := view.RenderString()
	if err != nil {
		return err
	}

	// For debug, print message
	//fmt.Printf("SENDING MAIL:\n%s", html)

	message.SetHTML(html)

	return sg.Send(message)
}

// SendOne sends email to ONE recipient only
func SendOne(recipient string, subject string, template string, context map[string]interface{}) error {
	return Send([]string{recipient}, subject, template, context)
}

// This is a dummy internal render context which doesn't provide any info - we explicitly set our info
// TODO: Perhaps find a more elegant solution?

// RenderContext provides an empty context for rendering mail, we have no preferred path
type renderContext struct {
}

// Path returns an empty path
func (m *renderContext) Path() string {
	return ""
}

// RenderContext returns a nil context
func (m *renderContext) RenderContext() map[string]interface{} {
	return nil
}
