package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCoop(root *goquery.Document) (models.RecipeSchema, error) {
	dateCreated, _ := root.Find("meta[name='creation_date']").Attr("content")

	name, _ := root.Find("meta[name='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, "|")
	if ok {
		name = strings.TrimSpace(before)
	}

	description, _ := root.Find("meta[name='og:description']").Attr("content")

	nodes := root.Find(".IngredientList-content")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("ol.List--orderedRecipe")
	var instructions []string
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	image, _ := root.Find("picture img").First().Attr("src")
	image = strings.TrimPrefix(image, "//")

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DateCreated:   dateCreated,
		DatePublished: dateCreated,
		Description:   models.Description{Value: strings.TrimSpace(description)},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(root.Find("span:contains('portioner')").Text())},
	}, nil
}
