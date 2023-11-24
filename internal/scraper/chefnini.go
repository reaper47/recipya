package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"strings"
)

func scrapeChefnini(root *goquery.Document) (models.RecipeSchema, error) {
	nodes := root.Find("meta[property='article:tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s, _ := sel.Attr("content")
		xk = append(xk, s)
	})
	keywords := strings.Join(extensions.Unique(xk), ",")

	categories, _ := root.Find("meta[property='article:section']").Attr("content")
	var category string
	if categories != "" {
		category = strings.TrimSpace(strings.Split(categories, ",")[0])
	}

	image, _ := root.Find("meta[property='og:image']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}

	description := root.Find("p[itemprop='description']").Text()
	description = strings.TrimSpace(description)

	nodes = root.Find("li[itemprop='ingredients']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] p")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		if s == "" {
			return
		}
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		Yield:         models.Yield{Value: findYield(root.Find("h3[itemprop='recipeYield']").Text())},
	}, nil
}
