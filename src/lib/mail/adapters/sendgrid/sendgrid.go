package sendgrid

import (
	"errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	m "github.com/fragmenta/fragmenta-cms/src/lib/mail"
)

// Service sends mail via sendgrid and conforms to mail.Service.
type Service struct {
	from   string
	secret string
}

// New returns a new sendgrid Service.
func New(f string, s string) *Service {
	return &Service{
		from:   f,
		secret: s,
	}
}

// Send the given message to recipients, using the context to render it
func (s *Service) Send(email *m.Email) error {

	if s.secret == "" {
		return errors.New("mail: invalid mail settings")
	}

	// Set the default from if required
	if email.ReplyTo == "" {
		email.ReplyTo = s.from
	}

	// Check if other fields are filled in on email
	if email.Invalid() {
		return errors.New("mail: attempt to send invalid email")
	}

	// Create a sendgrid message with the byzantine sendgrid API
	sendgridContent := mail.NewContent("text/html", email.Body)
	var sendgridRecipients []*mail.Email
	for _, r := range email.Recipients {
		// We could possibly split recipients on <> to get email (e.g. name<example@example.com>)
		// for now we assume they are just an email address
		sendgridRecipients = append(sendgridRecipients, mail.NewEmail("", r))
	}

	message := mail.NewV3Mail()
	message.Subject = email.Subject
	message.From = mail.NewEmail("", email.ReplyTo)
	if email.ReplyTo != "" {
		message.SetReplyTo(mail.NewEmail("", email.ReplyTo))
	}
	p := mail.NewPersonalization()
	p.AddTos(sendgridRecipients...)
	message.AddPersonalizations(p)
	message.AddContent(sendgridContent)

	request := sendgrid.GetRequest(s.secret, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(message)
	_, err := sendgrid.API(request)
	return err
}
