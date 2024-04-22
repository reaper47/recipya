package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"time"
)

func scrapeBongeats(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find(".recipe-ingredients li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".recipe-process li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	rs.Ingredients.Values = ingredients
	rs.Instructions.Values = instructions

	parsed, err := time.Parse("Jan 02, 2006", rs.DatePublished)
	if err == nil {
		rs.DatePublished = parsed.Format(time.DateOnly)
	}

	return rs, nil
}
