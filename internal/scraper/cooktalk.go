package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCooktalk(root *goquery.Document) (models.RecipeSchema, error) {
	datePublished, _ := root.Find("time.entry-date").Attr("datetime")

	nodes := root.Find("a[rel='category']")
	xc := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xc = append(xc, sel.Text())
	})

	image, _ := root.Find("img[itemprop='image']").Attr("src")

	description := root.Find("div[itemprop='description']").Text()
	description = strings.TrimSpace(strings.Trim(description, "\n"))

	nodes = root.Find("li[itemprop='ingredients']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		s = strings.TrimSpace(s)
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("p[itemprop='recipeInstructions']")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		s = strings.TrimSpace(s)
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: xc[0]},
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          root.Find(".page-title").Text(),
	}, nil
}
