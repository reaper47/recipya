package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeReisHunger(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find("span[ingredient-amount='']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.TrimSpace(s.Parent().Text())
	})

	var instructions []string
	root.Find("#zubereitung").Next().Find("div").Each(func(i int, s *goquery.Selection) {
		v := strings.TrimFunc(s.Text(), func(r rune) bool { return r == '\n' })
		if !strings.HasPrefix(v, "Schritt") {
			instructions = append(instructions, v)
		}
	})

	rs.Instructions = models.Instructions{Values: instructions}
	rs.Ingredients = models.Ingredients{Values: ingredients}
	return rs, nil
}
