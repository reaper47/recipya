package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeSaboresajinomoto(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	getIngredients(&rs, root.Find("div[itemprop='recipeIngredient'] ul li"))
	getInstructions(&rs, root.Find("div[itemprop='recipeInstructions'] ol li"))

	return rs, nil
}
