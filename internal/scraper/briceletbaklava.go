package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBriceletbaklava(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}

	image, _ := root.Find("meta[property='og:image']").Attr("content")

	content := root.Find(".ob-section-html")
	description := strings.Trim(content.First().Text(), "\n")
	description = strings.TrimSpace(description)

	nodes := root.Find(".Post-tags a")
	var xk []string
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s == "" {
			return
		}
		xk = append(xk, s)
	})
	keywords := strings.Join(xk, ",")

	nodes = content.Last().Find("p")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s == " " {
			return
		}
		ingredients = append(ingredients, strings.TrimSpace(s))
	})

	nodes = content.Last().Find("ul li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, strings.TrimSpace(s))
	})

	return models.RecipeSchema{
		AtType:       models.SchemaType{Value: "Recipe"},
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Keywords:     models.Keywords{Values: keywords},
		Name:         name,
	}, nil
}
