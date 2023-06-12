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

	yield := findYield(root.Find(".spritePortion").Siblings().First().Text())
	rs.Yield.Value = yield
	return rs, nil
}
