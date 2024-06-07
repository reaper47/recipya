package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeMindMegette(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	rs.Yield.Value = findYield(root.Find(".spritePortion").Siblings().First().Text())
	return rs, nil
}
