package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeFarmhousedelivery(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}

	description, _ := root.Find("meta[property='og:description']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	category := root.Find("a[rel='category tag']").First().Text()

	content := root.Find(".entry-content")
	var ingredients []string
	content.Find("ul li").Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	var instructions []string
	node := root.Find("p:contains('Instructions')")
	for {
		if node.Nodes == nil || goquery.NodeName(node) == "footer" {
			break
		}

		node = node.Next()
		s := strings.TrimSpace(node.Text())
		if s != "" {
			instructions = append(instructions, s)
		}
	}

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
	}, nil
}
