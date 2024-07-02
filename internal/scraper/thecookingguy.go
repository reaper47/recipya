package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"time"
)

func scrapeTheCookingGuy(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	date, err := time.Parse("Jan 02, 2006", rs.DatePublished)
	if err == nil {
		rs.DatePublished = date.Format(time.DateOnly)
	}

	getIngredients(&rs, root.Find("div.ingredients li"))
	getInstructions(&rs, root.Find("div.directions li"))

	node := root.Find("div.content_indent").Children().Last()
	if node != nil {
		rs.Yield.Value = findYield(node.Text())
	}

	return rs, nil
}
