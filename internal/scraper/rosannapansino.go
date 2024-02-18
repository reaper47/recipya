package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRosannapansino(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")

	description, _ := root.Find("meta[property='og:description']").Attr("content")
	before, _, found := strings.Cut(strings.TrimSpace(description), "\n")
	if found {
		description = strings.TrimSpace(before)
	}

	datePublished, _ := root.Find(".article__date time").First().Attr("datetime")

	content := root.Find(".recipe-content")

	nodes := content.Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = content.Find("ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		instructions = append(instructions, s)
	})

	image, _ := root.Find("div.article__body img").First().Attr("src")

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(root.Find("p:contains('Makes')").Text())},
	}, nil
}
