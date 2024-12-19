package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeReisHunger(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find("span[ingredient-amount='']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.TrimSpace(s.Parent().Text())
	})

	var instructions []models.HowToItem
	root.Find("#zubereitung").Next().Find("div").Each(func(_ int, s *goquery.Selection) {
		v := strings.TrimFunc(s.Text(), func(r rune) bool { return r == '\n' })
		if !strings.HasPrefix(v, "Schritt") {
			instructions = append(instructions, models.NewHowToStep(v))
		}
	})

	rs.Ingredients = &models.Ingredients{Values: ingredients}
	rs.Instructions = &models.Instructions{Values: slices.DeleteFunc(instructions, func(item models.HowToItem) bool {
		return strings.HasPrefix(item.Text, "Schritt")
	})}
	return rs, nil
}
