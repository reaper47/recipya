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

	if len(rs.Ingredients.Values) == 0 {
		nodes := root.Find("ul.wprm-recipe-ingredients").First().Find("li")
		rs.Ingredients.Values = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
		})
	}

	nodes := root.Find(".wprm-recipe-instructions-container").First().Find("ul.wprm-recipe-instructions li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})
	return rs, nil
}
