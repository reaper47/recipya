package regex

import "regexp"

// Email verifies whether an email address is valid.
var Email = regexp.MustCompile("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}$")

// Quantity detects quantities, i.e. 1ml, 1 ml, 1l and 1 l.
var Quantity = regexp.MustCompile(`\d+\s*(ml|mL|l|L)?(°[cCfF])?`)

// Regular expressions related to HTML.
var (
	// URL finds an HTTP or HTTPS address.
	URL = regexp.MustCompile(`^(https?://)?(?:[^@\n]+@)?([\da-z.-]+)\.([a-z.]{2,6})(:[0-9]+)?([/\w .-]*)*/?(?:[?&#][^\n]*)?$`)

	// Anchor is the regex to find all the content within an anchor tag.
	Anchor = regexp.MustCompile(`<a\s([\w =\\":;/.-])*\s*>`)

	// HourMinutes is the regex for the hour:minutes convention.
	HourMinutes = regexp.MustCompile(`^\d+:[0-5](\d?){1,2}?`)
)

// Sentences matches all sentences in a paragraph.
var Sentences = regexp.MustCompile(`(?:[-|\w,'\s%]*(?:\d*\.\d+%?)*[\s|\w]*)*[.|?!]`)

// Letters is the regular expression to search for letters in the text.
var Letters = regexp.MustCompile("[a-zA-Z]+")
