package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapePickupLimes(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if rs.DatePublished != "" {
		rs.DatePublished = ""
	}

	getInstructions(&rs, root.Find("ol").First().Find("li"))

	return rs, nil
}
