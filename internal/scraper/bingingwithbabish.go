package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBingingWithBabish(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " — ")
	if ok {
		name = strings.TrimSpace(before)
	}

	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	datePub, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	dateMod, _ := root.Find("meta[itemprop='dateModified']").Attr("content")

	nodes := root.Find("h3:contains('Ingredients')").First().Parent().Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("h1:contains('Method')").First().Parent().Find("ol li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DateModified:  dateMod,
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{1},
	}, nil
}
