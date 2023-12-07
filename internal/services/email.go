package services

import (
	"errors"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"jaytaylor.com/html2text"
	"strconv"
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

// RateLimits gets the SendGrid API's remaining and reset rate limits.
func (e *Email) RateLimits() (remaining int, resetUnix int64, err error) {
	req := sendgrid.GetRequest(app.Config.Email.SendGridAPIKey, "/v3/templates", "https://api.sendgrid.com")

	res, err := sendgrid.API(req)
	if err != nil {
		return -1, -1, err
	}

	xs, ok := res.Headers["X-Ratelimit-Remaining"]
	if !ok {
		return -1, -1, errors.New("cannot find the X-RateLimit-Remaining header")
	}

	rem, err := strconv.Atoi(xs[0])
	if err != nil {
		return -1, -1, err
	}

	xs, ok = res.Headers["X-Ratelimit-Reset"]
	if !ok {
		return -1, -1, errors.New("cannot find the X-RateLimit-Reset header")
	}

	reset, err := strconv.ParseInt(xs[0], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	return rem, reset, nil
}

// Send sends an email using the SendGrid API.
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
