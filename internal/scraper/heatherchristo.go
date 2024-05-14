package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHeatherChristo(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	prep, _ := root.Find("time[itemprop='prepTime']").Attr("datetime")
	cook, _ := root.Find("time[itemprop='cookTime']").Attr("datetime")

	nodes := root.Find(".ERSIngredients ul").First().Find("li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".ERSInstructions ol").First().Find("li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	return models.RecipeSchema{
		CookTime:        cook,
		Description:     models.Description{Value: description},
		Image:           models.Image{Value: image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            strings.TrimSpace(root.Find("div[itemprop='name']").First().Text()),
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        prep,
		Tools:           models.Tools{},
		Yield:           models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").First().Text())},
	}, nil
}
