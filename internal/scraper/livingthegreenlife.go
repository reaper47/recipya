package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeLivingTheGreenLife(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	var ingredients []string
	if len(ingredients) == 0 {
		nodes := root.Find("ul.wprm-recipe-ingredients").First().Find("li")
		ingredients = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			ingredients = append(ingredients, sel.Text())
		})
	}
	rs.Ingredients = &models.Ingredients{Values: ingredients}

	nodes := root.Find(".wprm-recipe-instructions-container").First().Find("ul.wprm-recipe-instructions li")
	instructions := make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(sel.Text()))
	})
	rs.Instructions = &models.Instructions{Values: instructions}

	return rs, nil
}
