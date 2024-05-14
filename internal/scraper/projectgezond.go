package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeProjectgezond(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " | ")
	if found {
		name = strings.TrimSpace(before)
	}
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	category, _ := root.Find("meta[property='article:section']").Attr("content")
	image, _ := root.Find(".wp-post-image").First().Attr("src")

	datePublished, _ := root.Find("meta[property='og:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")

	nodes := root.Find("h2").First().NextUntil("h2")
	ingredientNodes := nodes.Find("ul li")
	ingredients := make([]string, 0, ingredientNodes.Length())
	ingredientNodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = nodes.Next().NextUntil("h2")
	instructionNodes := nodes.Find("ul li")
	instructions := make([]string, 0, instructionNodes.Length())
	instructionNodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	var cal string
	node := root.Find("strong:contains('Kcal')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		cal = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var protein string
	node = root.Find("strong:contains('Eiwit')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		protein = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var carbs string
	node = root.Find("strong:contains('Koolhydraten')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		carbs = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var fat string
	node = root.Find("strong:contains('Vet')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		fat = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var fiber string
	node = root.Find("strong:contains('Vezels')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		fiber = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	return models.RecipeSchema{
		Category:      models.Category{Value: category},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		NutritionSchema: models.NutritionSchema{
			Calories:      cal + " kcal",
			Carbohydrates: carbs,
			Fat:           fat,
			Fiber:         fiber,
			Protein:       protein,
		},
	}, nil
}
