package email_test

import (
	"os"
	"testing"

	"github.com/reaper47/recipya/internal/email"
)

func TestEmail(t *testing.T) {
	t.Run("email is valid", func(t *testing.T) {
		setEnv()
		defer unsetEnv()

		if !email.IsValid() {
			t.Error("got invalid email config when it should have been")
		}
	})

	t.Run("need to refresh environment variables", func(t *testing.T) {
		email.Reload()

		if email.IsValid() {
			t.Error("got valid email config when it should not have been")
		}
	})
}

func setEnv() {
	_ = os.Setenv("MAIL_HOST", "host")
	_ = os.Setenv("MAIL_TO", "to")
	_ = os.Setenv("MAIL_PASSWORD", "password")
	_ = os.Setenv("MAIL_PORT", "1234")
}

func unsetEnv() {
	_ = os.Unsetenv("MAIL_TO")
	_ = os.Unsetenv("MAIL_HOST")
	_ = os.Unsetenv("MAIL_PASSWORD")
	_ = os.Unsetenv("MAIL_PORT")
}
