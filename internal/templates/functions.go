package templates

import (
	"fmt"
	"net/url"
	"text/template"
	"time"
)

var fm = template.FuncMap{
	"fmtDuration": func(d time.Duration) string {
		d = d.Round(time.Minute)
		h := d / time.Hour
		d -= h * time.Hour
		m := d / time.Minute
		d -= m * time.Minute
		return fmt.Sprintf("%dh%02d", h, m)
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
