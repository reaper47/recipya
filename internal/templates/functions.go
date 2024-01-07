package templates

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/url"
	"strings"
	"time"
)

func formatDuration(d time.Duration, isDatetime bool) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute

	d -= m * time.Minute
	if isDatetime {
		if h > 0 {
			return fmt.Sprintf("PT%dH%02dM", h, m)
		}
		return fmt.Sprintf("PT%02dM", m)
	}

	return fmt.Sprintf("%dh%02dm", h, m)
}

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		if strings.HasPrefix(s, "http") {
			log.Printf("isURL.ParseRequestURI error: %q", err)
		}
		return false
	}

	u, err := url.Parse(s)
	if u == nil {
		log.Printf("parsed URL %q is nil", s)
		return false
	}
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isUUIDValid(u uuid.UUID) bool {
	return u != uuid.Nil
}
