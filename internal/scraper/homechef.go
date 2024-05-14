package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHomechef(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")
	yield, _ := root.Find("meta[itemprop='recipeYield']").Attr("content")

	image, _ := root.Find("div img").First().Attr("data-srcset")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		s = strings.TrimSpace(strings.TrimPrefix(s, "Info"))
		if s != "" {
			ingredients = append(ingredients, s)
		}
	})

	nodes = root.Find("li[itemprop='itemListElement']")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		if s != "" {
			instructions = append(instructions, s)
		}
	})

	return models.RecipeSchema{
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find("h1").First().Text(),
		Yield:        models.Yield{Value: findYield(yield)},
	}, nil
}
