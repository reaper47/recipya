package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeBongeats(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find(".recipe-ingredients li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".recipe-process li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	rs.Ingredients.Values = ingredients
	rs.Instructions.Values = instructions

	return rs, nil
}
