package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeGesundAktiv(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")

	nodes := root.Find("div.field--name-field-zutaten .field--item")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("div.field--name-field-content-element p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Description:  models.Description{Value: description},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         strings.TrimSpace(root.Find("h1.page-header").Text()),
	}, nil
}
