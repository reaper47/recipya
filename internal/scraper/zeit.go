package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeZeit(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")
	dateModified, _ := root.Find("meta[name='last-modified']").Attr("content")
	datePublished, _ := root.Find("meta[name='date']").Attr("content")
	keywords, _ := root.Find("meta[name='keywords']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")

	meta := root.Find(".recipe-list-meta")

	node := meta.Find("title:contains('Portionen')")
	var yield string
	if node.Parent() != nil && node.Parent().Parent() != nil {
		yield = node.Parent().Parent().Text()
	}

	nodes := root.Find(".recipe-list-collection__list-item")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".paragraph.article__item")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(yield)},
	}, nil
}
