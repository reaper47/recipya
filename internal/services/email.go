package services

import (
	"errors"
	"github.com/k3a/html2text"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/wneessen/go-mail"
	"log/slog"
)

type request struct {
	from    string
	to      string
	subject string
	body    string
}

type emailItem struct {
	to       string
	template templates.EmailTemplate
	data     any
}

// NewEmailService creates a new Email service.
func NewEmailService() *Email {
	return &Email{
		queue: make([]emailItem, 0),
	}
}

// Email is the entity that manages the email client.
type Email struct {
	queue []emailItem
}

// Queue adds an unsent email to the queue.
func (e *Email) Queue(to string, template templates.EmailTemplate, data any) {
	e.queue = append(e.queue, emailItem{
		to:       to,
		template: template,
		data:     data,
	})
}

// RateLimits gets the SMTP server's remaining and reset rate limits.
func (e *Email) RateLimits() (remaining int, resetUnix int64, err error) {
	return 1000, 1000, nil
}

// Send sends an email using SMTP.
func (e *Email) Send(to string, template templates.EmailTemplate, data any) error {
	r := &request{
		from:    app.Config.Email.From,
		to:      to,
		subject: template.Subject(),
	}

	buf := templates.RenderEmail(template.String(), data)
	r.body = buf

	return r.sendMail()
}

// SendQueue sends emails in the queue until the rate limit has been reached.
func (e *Email) SendQueue() (sent, remaining int, err error) {
	if len(e.queue) == 0 {
		return 0, 0, errors.New("no emails in queue")
	}

	remaining, _, err = e.RateLimits()
	if err != nil {
		return sent, remaining, err
	}

	if remaining > len(e.queue) {
		remaining = len(e.queue)
	}

	count := 0
	for i := remaining; i > 0; i-- {
		el := e.queue[0]

		err = e.Send(el.to, el.template, el.data)
		if err != nil {
			return count, len(e.queue), err
		}

		e.queue = e.queue[1:]
		count++
	}

	return count, len(e.queue), nil
}

func (r *request) sendMail() error {
	go func() {
		m := mail.NewMsg()
		if err := m.From(r.from); err != nil {
			slog.Error("SendEmail failed to set from", "error", err)
			return
		}
		if err := m.To(r.to); err != nil {
			slog.Error("SendEmail failed to set to", "error", err)
			return
		}
		m.Subject(r.subject)
		m.SetBodyString(mail.TypeTextPlain, html2text.HTML2Text(r.body))
		m.AddAlternativeString(mail.TypeTextHTML, r.body)

		client, err := mail.NewClient(
			app.Config.Email.Host,
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(app.Config.Email.Username),
			mail.WithPassword(app.Config.Email.Password),
		)
		if err != nil {
			slog.Error("SendEmail failed to create client", "error", err)
			return
		}

		err = client.DialAndSend(m)
		if err != nil {
			slog.Error("SendEmail failed to send email", "error", err)
			return
		}
	}()

	return nil
}
