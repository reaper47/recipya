package regex

import "regexp"

// Email is the regex used to verify whether an email address is valid.
var Email = regexp.MustCompile(
	"^[\\w.!#$%&'*+/=?^_`{|}~-]+@[\\w](?:[\\w-]{0,61}[\\w])?(?:\\.[\\w](?:[\\w-]{0,61}[\\w])?)*$",
)

// Regular expressions related to HTML.
var (
	// Anchor is the regex to find all the content within an anchor tag.
	Anchor = regexp.MustCompile(`<a([\w =\\":;\/.-])*\s*>`)

	// HourMinutes is the regex for the hour:minutes convention.
	HourMinutes = regexp.MustCompile(`\d+:[0-5](\d?){1,2}?`)

	// ImageSrc is the regex to find a URL within the source attribute of an image tag
	ImageSrc = regexp.MustCompile(`https://[\w./?=%-_&;]+`)
)

// Quantity is the regex for detecting quantities, i.e. 1ml, 1 ml, 1l and 1 l.
var Quantity = regexp.MustCompile(`([\d]\s*m?[l]{1})?(\d*°[cf]{1})?`)
