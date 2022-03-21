package logger

import (
	"log"
	"strings"
)

// Sanitize saniztize a log.Println message by removing new
// lines and carriage returns.
func Sanitize(message string, args ...string) {
	if len(args) > 0 {
		var sanitized []string
		for _, arg := range args {
			escaped := strings.Replace(arg, "\n", "", -1)
			escaped = strings.Replace(escaped, "\r", "", -1)
			sanitized = append(sanitized, escaped)
		}
		log.Println(message, sanitized)
	} else {
		log.Println(message)
	}
}
