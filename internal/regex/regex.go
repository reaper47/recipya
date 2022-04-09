package regex

import "regexp"

// Email is the regex used to verify whether an email address is valid.
var Email = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Regular expressions related to HTML.
var (
	// Anchor is the regex to find all the content within an anchor tag.
	Anchor = regexp.MustCompile(`<a([A-Za-z =\\":\/.-])*>`)

	// HourMinutes is the regex for the hour:minutes convention.
	HourMinutes = regexp.MustCompile(`[0-9]+:[0-9]+`)

	// ImageSrc is the regex to find a URL within the source attribute of an image tag
	ImageSrc = regexp.MustCompile("https://[a-z./0-9?=%A-Z-_]+")
)
