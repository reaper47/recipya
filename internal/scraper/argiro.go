package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeArgiro(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	parts := strings.Split(rs.DatePublished, "-")
	if len(parts) == 3 {
		slices.Reverse(parts)

		if len(parts[0]) == 1 {
			parts[0] = "0" + parts[0]
		}

		if len(parts[1]) == 1 {
			parts[1] = "0" + parts[1]
		}

		rs.DatePublished = strings.Join(parts, "-")
	}

	return rs, nil
}
