package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeDrinkoteket(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Last().Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	category, _ := root.Find("meta[property='article:section']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Last().Attr("content")

	nodes := root.Find("ul.ingredients li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("#recipe-utrustning .rbs-img-content")
	tools := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		tools = append(tools, strings.TrimSpace(sel.Text()))
	})

	datePub, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	prep, _ := root.Find("meta[itemprop='prepTime']").Attr("content")
	cook, _ := root.Find("meta[itemprop='cookTime']").Attr("content")

	return models.RecipeSchema{
		Category:      models.Category{Value: category},
		CookTime:      cook,
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Tools:         models.Tools{Values: tools},
		Yield:         models.Yield{Value: 1},
	}, nil
}
