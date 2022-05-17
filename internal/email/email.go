package email

import (
	"log"
	"net/smtp"
	"os"
	"sync"
)

var (
	email config
	once  sync.Once
)

type config struct {
	Auth smtp.Auth
	Addr string
	To   string
}

// Send sends an email to one recipient.
func Send(from, msg string) {
	once.Do(Reload)
	msg = "To: " + email.To + "\r\n" + msg
	err := smtp.SendMail(email.Addr, email.Auth, from, []string{email.To}, []byte(msg))
	if err != nil {
		log.Printf("could not send email: %s", err)
	}
}

// IsValid verifies whether the Config struct is configured. If it is not, then
// the user did not set the MAIL_* variables under the .env file.
func IsValid() bool {
	once.Do(Reload)
	return email.Addr != "" && email.To != ""
}

// Reload refreshes the email's configuration from the environment variables.
func Reload() {
	host := os.Getenv("MAIL_HOST")
	to := os.Getenv("MAIL_TO")

	email = config{
		Auth: smtp.PlainAuth("", to, os.Getenv("MAIL_PASSWORD"), host),
		Addr: host + ":" + os.Getenv("MAIL_PORT"),
		To:   to,
	}
}
