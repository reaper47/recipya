package templates

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
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
			slog.Info("Failed to parse request URI", "uri", s, "error", err)
		}
		return false
	}

	u, err := url.Parse(s)
	if u == nil {
		slog.Error("Parsed URL is nil", "url", s, "error", err)
		return false
	}
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isUUIDsValid(xu []uuid.UUID) []bool {
	xb := make([]bool, 0, len(xu))
	for _, u := range xu {
		xb = append(xb, u != uuid.Nil)
	}
	return xb
}
