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
	return regex.Digit.ReplaceAllStringFunc(s, func(s string) string {
		sum := SumString(s)
		return FloatToString(sum*scale, "%f")
	})
}

// SumString sums consecutive numbers in a string.
func SumString(s string) float64 {
	sum := 0.

	matches := regex.DimensionPattern.FindAllStringSubmatch(s, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			l, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				continue
			}

			r, err := strconv.ParseFloat(match[2], 64)
			if err != nil {
				continue
			}

			sum += l * r
		}

		s = regex.DimensionPattern.ReplaceAllString(s, "")
	}

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
