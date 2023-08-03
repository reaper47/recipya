package regex

import "regexp"

// Email verifies whether an email address is valid.
var Email = regexp.MustCompile("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}$")

// Quantity detects quantities, i.e. 1ml, 1 ml, 1l and 1 l.
var Quantity = regexp.MustCompile(`\d+\s*((ml|mL|l|L)(°[cCfF])?|°[cCfF])`)

// Regular expressions related to HTML.
var (
	// URL finds an HTTP or HTTPS address.
	URL = regexp.MustCompile(`^(https?://)?(?:[^@\n]+@)?([\da-z.-]+)\.([a-z.]{2,6})(:[0-9]+)?([/\w .-]*)*/?(?:[?&#][^\n]*)?$`)

	// Anchor is the regex to find all the content within an anchor tag.
	Anchor = regexp.MustCompile(`<a\s([\w =\\":;/.-])*\s*>`)

	// HourMinutes is the regex for the hour:minutes convention.
	HourMinutes = regexp.MustCompile(`^\d+:[0-5](\d?){1,2}?`)
)

// Letters is the regular expression to search for letters in the text.
var Letters = regexp.MustCompile("[a-zA-Z]+")

// Unit is the regular expression to search for
var Unit = regexp.MustCompile(`(?i)((?:\d*\.?\d+\s*to\s*)?(?:\d*\s*\d+/)?(?:\d+-\d*/?)?\d*\.?\d+)-?\s*(centimeters?|centimetres?|cm\b|cups?|deciliters?|decilitres?|dl\b|feet|foot|ft\.?|′|fluid\s*ounces|fl\.?\s*oz\.*|fluid\s*oz\.?|gallons?|gals?\b|milliliters?|millilitres?|ml\b|millimeters?|millimetres?|mm\b|grams?|grammes?|g\b|inches?|inch|in\b|["”]|kilograms?|kilogrammes?|kg|milligrams?|milligrammes?|mg\b|meters?|metres?|m\b|ounces?|oz\.?|pints?|fl\.?\s*pt\.?|pt\.?|pounds?|lbs?\.?|lb\.?|#|quarts?|fl\.?\s*qt\.?|qt\.?|liters?|litres?|l\b|tablespoons?|tbsp\.?\w*|teaspoons?|tsp\.?\w*|yards?|degrees?\s*celsius|degrees?\s*c|celsius|°?c\b|degrees?\s*fahrenheit|degrees?\s*f|fahrenheit|°?f\b)`)
