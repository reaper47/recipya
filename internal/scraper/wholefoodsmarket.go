package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeWholefoodsmarket(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[itemprop='headline']").Attr("content")
	dateModified, _ := root.Find("meta[itemprop='dateModified']").Attr("content")
	datePublished, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	image, _ := root.Find("meta[itemprop='image']").Attr("content")
	description, _ := root.Find("meta[itemprop='description']").Attr("content")

	p := root.Find(".image-subtitle p").Last().Text()
	var yield string
	for _, s := range strings.Split(p, "|") {
		if strings.Contains(strings.ToLower(s), "serves") {
			yield = s
		}
	}

	nodes := root.Find("h4:contains('Ingredients')").Parent().Find("p")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("h4:contains('Method')").ParentsUntil(".sqs-col-6").Last().Parent().Find("p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		if s != "" {
			instructions = append(instructions, s)
		}

	})

	return models.RecipeSchema{
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(yield)},
	}, nil
}
