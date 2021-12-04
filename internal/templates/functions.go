package templates

import (
	"fmt"
	"net/url"
	"text/template"
	"time"

	"github.com/google/uuid"
)

var fm = template.FuncMap{
	"durationToInput": func(d time.Duration) string {
		h, m, d := getHourMin(d)
		s := d / time.Second
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	},
	"fmtDuration": func(d time.Duration) string {
		h, m, d := getHourMin(d)
		d -= m * time.Minute
		return fmt.Sprintf("%dh%02d", h, m)
	},
	"inc": func(i int) int {
		return i + 1
	},
	"isUrl": func(s string) bool {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			return false
		}

		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return false
		}
		return true
	},
	"isUuidValid": func(u uuid.UUID) bool {
		return u != uuid.UUID{}
	},
	"substring": func(s string, endIndex int) string {
		if len(s) <= endIndex {
			return s
		}
		return s[:endIndex] + "..."
	},
}

func getHourMin(d time.Duration) (time.Duration, time.Duration, time.Duration) {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	return h, m, d
}
