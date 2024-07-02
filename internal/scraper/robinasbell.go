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

	img, _ := root.Find(".wp-caption img").First().Attr("data-src")
	rs.Image = &models.Image{Value: img}

	return rs, nil
}
