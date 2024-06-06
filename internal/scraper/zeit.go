package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeZeit(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.DateModified, _ = root.Find("meta[name='last-modified']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[name='date']").Attr("content")
	keywords, _ := root.Find("meta[name='keywords']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")

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
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		instructions = append(instructions, models.NewHowToStep(s))
	})

	return models.RecipeSchema{
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   &models.Description{Value: description},
		Keywords:      &models.Keywords{Values: keywords},
		Image:         &models.Image{Value: image},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
		Name:          name,
		Yield:         &models.Yield{Value: findYield(yield)},
	}, nil
}
