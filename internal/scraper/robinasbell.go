package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeRobinasBell(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Image = &models.Image{Value: root.Find(".wp-caption img").First().AttrOr("data-src", "")}
	return rs, nil
}
