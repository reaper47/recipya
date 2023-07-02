package templates_test

import (
	"github.com/reaper47/recipya/internal/templates"
	"testing"
)

var emailTemplates = []templates.EmailTemplate{
	templates.EmailErrorAdmin,
	templates.EmailForgotPassword,
	templates.EmailIntro,
	templates.EmailRequestWebsite,
}

func TestEmailTemplate_String(t *testing.T) {
	want := []string{
		"error-admin.mjml",
		"forgot-password.mjml",
		"intro.mjml",
		"request-website.mjml",
	}
	for i, template := range emailTemplates {
		if got := template.String(); got != want[i] {
			t.Fatalf("got %q but want %q", got, want[i])
		}
	}
}

func TestEmailTemplate_Subject(t *testing.T) {
	want := []string{
		"Recipya Error",
		"Forgot Password",
		"Confirm Account",
		"Request Website",
	}
	for i, template := range emailTemplates {
		if got := template.Subject(); got != want[i] {
			t.Fatalf("got %q but want %q", got, want[i])
		}
	}
}
