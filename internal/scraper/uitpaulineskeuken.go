package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeUitpaulineskeuken(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	before, _, found := strings.Cut(rs.Category.Value, ",")
	if found {
		rs.Category.Value = strings.TrimSpace(before)
	}

	return rs, nil
}
