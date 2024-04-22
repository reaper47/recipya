package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeKingArthurBaking(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if rs.DatePublished != "" {
		before, _, _ := strings.Cut(rs.DatePublished, "at")

		t, err := time.Parse("January 2, 2006", strings.TrimSpace(before))
		if err == nil {
			rs.DatePublished = t.Format(time.DateOnly)
		}
	}

	return rs, nil
}
