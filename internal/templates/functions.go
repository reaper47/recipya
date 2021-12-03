package templates

import (
	"fmt"
	"net/url"
	"text/template"
	"time"
)

var fm = template.FuncMap{
	"fmtDuration": func(d time.Duration) string {
		h, m, d := getHourMin(d)
		d -= m * time.Minute
		return fmt.Sprintf("%dh%02d", h, m)
	},
	"durationToInput": func(d time.Duration) string {
		h, m, d := getHourMin(d)
		s := d / time.Second
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
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
}

func getHourMin(d time.Duration) (time.Duration, time.Duration, time.Duration) {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	return h, m, d
}
