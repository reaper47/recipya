package models

import (
	"strconv"
	"strings"
)

// Websites is the type for a slice of Website.
type Websites []Website

// TableHTML renders the websites to a string of table rows.
func (w Websites) TableHTML() string {
	var sb strings.Builder
	for _, website := range w {
		tr := `<tr class="border px-8 py-2">`
		data1 := `<td class="border px-8 py-2">` + strconv.FormatInt(website.ID, 10) + "</td>"
		link := `<a class="underline" href="` + website.URL + `" target="_blank">` + website.Host + "</a>"
		data2 := `<td class="border px-8 py-2">` + link + "</td>"
		sb.WriteString(tr + data1 + data2 + "</tr>")
	}

	if sb.Len() == 0 {
		tr := `<tr class="border px-8 py-2">`
		data1 := `<td class="border px-8 py-2">-1</td>`
		data2 := `<td class="border px-8 py-2">No result</td>`
		sb.WriteString(tr + data1 + data2 + "</tr>")
	}

	return sb.String()
}

// Website represents a website with its hostname and URL.
type Website struct {
	ID   int64
	Host string
	URL  string
}
