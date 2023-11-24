package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCostco(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")

	keywordsRaw, _ := root.Find("meta[name='description']").Attr("content")
	var keywords strings.Builder
	for _, s := range strings.Split(keywordsRaw, ",") {
		if s != "" {
			keywords.WriteString(s)
		}
	}

	h1 := root.Find("h1").Last()
	div := h1.Parent()
	name := h1.Text()

	image, _ := div.Prev().Find("img").Attr("src")

	nodes := div.Find("ul li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		ingredients = append(ingredients, s)
	})

	nodes = div.Find("p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			instructions = append(instructions, s)
		}
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		Description:  models.Description{Value: description},
		Keywords:     models.Keywords{Values: keywords.String()},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         name,
	}, nil
}
