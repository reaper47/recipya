package extensions

import (
	"fmt"
	"strings"
)

// FloatToString converts a float to a string. Trailing zeroes will be trimmed.
// The decimal will be trimmed if no trailing zeroes are present.
func FloatToString(number float64, format string) string {
	formatted := fmt.Sprintf(format, number)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return formatted
}
