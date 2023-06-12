package services

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"jaytaylor.com/html2text"
	"log"
)

type request struct {
	from    string
	to      string
	subject string
	body    string
}

// NewEmailService creates a new Email service.
func NewEmailService() *Email {
	return &Email{}
}

// Email is the entity that manages the email client.
type Email struct{}

// Send sends an email using the SendGrid API.
func (e *Email) Send(to string, template templates.EmailTemplate, data any) {
	go func() {
		r := &request{
			from:    app.Config.Email.From,
			to:      to,
			subject: template.Subject(),
		}

		buf := templates.RenderEmail(template.String(), data)
		r.body = buf

		if err := r.sendMail(); err != nil {
			log.Printf("error sending %s email to %s: %q", template, to, err)
		}
	}()

}

func (r *request) sendMail() error {
	text, err := html2text.FromString(r.body, html2text.Options{TextOnly: false})
	if err != nil {
		return err
	}

	client := sendgrid.NewSendClient(app.Config.Email.SendGridAPIKey)

	from := mail.NewEmail("Recipya", r.from)
	to := mail.NewEmail(r.subject, r.to)
	_, err = client.Send(mail.NewSingleEmail(from, r.subject, to, text, r.body))
	return err
}
