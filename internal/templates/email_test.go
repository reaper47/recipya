package templates_test

import (
	"github.com/reaper47/recipya/internal/templates"
	"testing"
)

var emailTemplates = []templates.EmailTemplate{
	templates.EmailErrorAdmin,
	templates.EmailForgotPassword,
	templates.EmailIntro,
}

func TestEmailTemplate_String(t *testing.T) {
	want := []string{
		"error-admin.gohtml",
		"forgot-password.gohtml",
		"intro.gohtml",
	}
	for i, template := range emailTemplates {
		got := template.String()
		if got != want[i] {
			t.Fatalf("got %q but want %q", got, want[i])
		}
	}
}

func TestEmailTemplate_Subject(t *testing.T) {
	want := []string{
		"Recipya Error",
		"Forgot Password",
		"Confirm Account",
	}
	for i, template := range emailTemplates {
		got := template.Subject()
		if got != want[i] {
			t.Fatalf("got %q but want %q", got, want[i])
		}
	}
}
