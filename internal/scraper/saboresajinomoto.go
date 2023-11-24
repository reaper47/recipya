package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeSaboresajinomoto(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find("div[itemprop='recipeIngredient'] ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})
	rs.Ingredients.Values = ingredients

	nodes = root.Find("div[itemprop='recipeInstructions'] ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})
	rs.Instructions.Values = instructions

	return rs, nil
}
