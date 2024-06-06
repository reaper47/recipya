package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeYemek(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(s.Text()))
	})

	nodes = root.Find("p[itemprop='recipeInstructions']")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(strings.TrimSpace(s.Text())))
	})

	return models.RecipeSchema{
		Image:        &models.Image{Value: image},
		Ingredients:  &models.Ingredients{Values: ingredients},
		Instructions: &models.Instructions{Values: instructions},
		Name:         name,
	}, nil
}
