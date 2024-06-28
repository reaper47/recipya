package scraper

import (
	"github.com/PuerkitoBio/goquery"
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

func getItempropContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[itemprop='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}

func getNameContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[name='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}

func getPropertyContent(doc *goquery.Document, name string) string {
	s, _ := doc.Find("meta[property='" + name + "']").Attr("content")
	return strings.TrimSpace(s)
}
