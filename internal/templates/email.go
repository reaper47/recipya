package templates

import (
	"bytes"
)

// EmailTemplate represents the name of a .mjml email template.
type EmailTemplate string

// These constants associate an EmailTemplate with its MJML file.
const (
	EmailErrorAdmin     EmailTemplate = "error-admin.mjml"
	EmailForgotPassword EmailTemplate = "forgot-password.mjml"
	EmailIntro          EmailTemplate = "intro.mjml"
)

// String represents the email template as a string, being the file name.
func (e EmailTemplate) String() string {
	return string(e)
}

// Subject returns the subject of the email according to the type of email being sent.
func (e EmailTemplate) Subject() string {
	switch e {
	case EmailErrorAdmin:
		return "Recipya Error"
	case EmailForgotPassword:
		return "Forgot Password"
	case EmailIntro:
		return "Confirm Account"
	default:
		return ""
	}
}

// EmailData holds data for email templates.
type EmailData struct {
	Text     string // Text is the text for the email.
	Token    string // Token is used to store JWT tokens.
	UserName string // UserName is the name of the user.
	URL      string // URL is the url of the website.
}

// RenderEmail is a wrapper for template.ExecuteTemplate on email templates.
func RenderEmail(name string, data any) string {
	tmpl, ok := templatesEmail[name]
	if !ok {
		return ""
	}

	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, data)
	return buf.String()
}
