package scraper

import (
	"strconv"
	"strings"
)

func findYield(s string) int16 {
	parts := strings.Split(s, " ")
	for _, part := range parts {
		i, err := strconv.ParseInt(part, 10, 16)
		if err == nil {
			return int16(i)
		}
	}
	return 0
}
