package extensions

import (
	"fmt"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
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

// ScaleString scales the numbers in the string in-place. The string may contain fractions.
func ScaleString(s string, scale float64) string {
	sum := SumString(s)
	var sb strings.Builder
	if sum > 0 {
		sb.WriteString(FloatToString(sum*scale, "%f"))
	}

	start := 0
	matches := regex.Digit.FindAllStringIndex(s, -1)
	for _, i := range matches {
		sb.WriteString(s[start:i[0]])
		start = i[1]
	}
	sb.WriteString(s[start:])
	return strings.Join(strings.Fields(sb.String()), " ")
}

// SumString sums consecutive numbers in a string.
func SumString(s string) float64 {
	sum := 0.
	for _, v := range strings.Split(s, " ") {
		if v == "" {
			continue
		}

		index := strings.Index(v, "/")
		if index != -1 {
			numerator, err := strconv.ParseFloat(v[:index], 64)
			if err != nil {
				continue
			}

			denominator, err := strconv.ParseFloat(v[index+1:], 64)
			if err != nil {
				continue
			}
			sum += numerator / denominator
		} else {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				if sum == 0. {
					continue
				}
				return sum
			}
			sum += f
		}
	}
	return sum
}
