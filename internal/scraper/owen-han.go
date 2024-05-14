package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeOwenhan(root *goquery.Document) (models.RecipeSchema, error) {
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	datePublished, _ := root.Find("meta[property='datePublished']").Attr("content")
	dateModified, _ := root.Find("meta[property='dateModified']").Attr("content")
	name, _ := root.Find("meta[itemprop='headline']").Attr("content")
	description, _ := root.Find("meta[itemprop='description']").Attr("content")

	content := root.Find("h4:contains('INGREDIENTS')").Parent()

	nodes := content.Find("ul p")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = content.Find("ol p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		Category:      models.Category{Value: root.Find(".blog-item-category").First().Text()},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
	}, nil
}
