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
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})
	rs.Ingredients.Values = ingredients

	nodes = root.Find("div[itemprop='recipeInstructions'] ol li")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(sel.Text()))
	})
	rs.Instructions.Values = instructions

	return rs, nil
}
