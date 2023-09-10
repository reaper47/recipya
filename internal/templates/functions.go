package templates

import (
	"fmt"
	"github.com/google/uuid"
	"net/url"
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
		return false
	}

	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isUUIDValid(u uuid.UUID) bool {
	return u != uuid.Nil
}
