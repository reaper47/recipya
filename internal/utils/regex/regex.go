package regex

import "regexp"

// BeginsWithWord matches a word at the beginning of a text.
var BeginsWithWord = regexp.MustCompile(`(?i)^[a-z]+[^\d]`)

// Decimal matches a decimal number.
var Decimal = regexp.MustCompile(`\d?\.\d+\b`)

// Digit matches digits in the text.
var Digit = regexp.MustCompile(`(\b\d+\s+\d+/\d+\b)|(\d+\.?/?\d*)`)

// DimensionPattern matches patterns representing dimensions.
var DimensionPattern = regexp.MustCompile(`(\d+)\s*x\s*(\d+).`)

// Email verifies whether an email address is valid.
var Email = regexp.MustCompile(`^[\w-.]+@([\w-]+\.)+[\w-]{2,4}$`)

// Quantity detects quantities, i.e. 1ml, 1 ml, 1l and 1 l.
var Quantity = regexp.MustCompile(`(?i)\d+\s*((ml|l\b)(°[cf])?|°[cf])`)

// Regular expressions related to HTML.
var (
	// URL matches an HTTP or HTTPS address.
	URL = regexp.MustCompile(`^(https?://)?(?:[^@\n]+@)?([\da-z.-]+)\.([a-z.]{2,6})(:[0-9]+)?([/\w .-]*)*/?(?:[?&#][^\n]*)?$`)

	// Anchor matches content within an anchor tag.
	Anchor = regexp.MustCompile(`<a\s([\w =\\":;/.-])*\s*>`)

	// HourMinutes matches the hour:minutes convention.
	HourMinutes = regexp.MustCompile(`^\d+:[0-5](\d?){1,2}?`)
)

// Letters matches letters in the text.
var Letters = regexp.MustCompile("[a-zA-Z]+")

// RangePattern matches numerical ranges.
var RangePattern = regexp.MustCompile(`(\d+(?:/\d+)?)\s*-\s*(\d+(?:/\d+)?)`)

// Time matches time, such as 1h30min.
var Time = regexp.MustCompile(`(?i)(\d+\s?h\s*)?(\d+\s?(?:m\b|min|minute|minutter|minuten|timer?)s?\b)|(\d+\s?h\s*)(\d+\s?mins?\b)?|(\d+\s?-\s?\d+\s*timer)`)

// Unit matches a unit.
var Unit = regexp.MustCompile(`(?i)((?:\d*\.?\d+\s*to\s*)?(?:\d*\s*\d+/)?(?:\d+-\d*/?)?\d*\.?\d+)-?\s*(centimeters?|centimetres?|cm\b|cups?|deciliters?|decilitres?|dl\b|feet|foot|ft\.?\b|′|fluid\s*ounces|fl\.?\s*oz\.*|fluid\s*oz\.?|gallons?|gals?\b|milliliters?|millilitres?|ml\b|millimeters?|millimetres?|mm\b|grams?|grammes?|\d*g\b|inches?|inch|in\b|["”]|kilograms?|kilogrammes?|kg|milligrams?|milligrammes?|mg\b|meters?|metres?|m\b|ounces?|oz\.?|pints?|fl\.?\s*pt\.?|pt\.?|pounds?|lbs?\.?\b|lb\.?\b|#|quarts?|fl\.?\s*qt\.?|qt\.?\b|liters?|litres?|l\b|tablespoons?|ss|tbsp\.?\w*|teaspoons?|ts\w?\.?|tsp\.?\w*|yards?|degrees?\s*celsius|degrees?\s*c|celsius|°?\s?c\b|degrees?\s*fahrenheit|degrees?\s*f|fahrenheit|°?\s?f\b)`)

// UnitImperial matches an imperial unit.
var UnitImperial = regexp.MustCompile(`(?i)(cups?|feet|foot|ft\.?\b|′|fluid\s*ounces|fl\.?\s*oz\.*|fluid\s*oz\.?|gallons?|gals?\b|inches?|inch|\d\s?in\b|["”]|ounces?|oz\.?|pints?|fl\.?\s*pt\.?|pt\.?\b|pounds?|lbs?\.?\b|lb\.?\b|#|quarts?|fl\.?\s*qt\.?|qt\.?\b|tablespoons?|tbsp\.?\w*|teaspoons?|tsp\.?\w*|yards?|degrees?\s*fahrenheit|degrees?\s*f|fahrenheit|\b°?f\b)`)

// UnitMetric matches a metric unit.
var UnitMetric = regexp.MustCompile(`(?i)(centimeters?|centimetres?|cm\b|deciliters?|decilitres?|dl\b|millimeters?|millimetres?|mm\b|grams?|grammes?|\b\d*g\b|kilograms?|kilogrammes?|kg|milligrams?|milligrammes?|mg\b|meters?|metres?|\b\d*m\b|milliliters?|millilitres?|ml\b|liters?|litres?|\b\d*l\b|degrees?\s*celsius|degrees?\s*c|celsius|\b°?c\b)`)

// WildcardURL matches a Recipya URL with wildcards.
var WildcardURL = regexp.MustCompile(`/(cookbooks|recipes|download)/\d+(/\w+(/\d+)?)?$`)
