package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrouprecipes(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")
	image, _ := root.Find(".photos img").First().Attr("src")
	cook, _ := root.Find(".cooktime .value-title").Attr("title")

	var keywords strings.Builder
	root.Find(".tags_text li").Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		keywords.WriteString(",")
		keywords.WriteString(s)
	})

	nodes := root.Find(".ingredients li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		before, _, ok := strings.Cut(s, "\t")
		if ok {
			s = strings.TrimSpace(before)
		}
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".instructions li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		before, after, ok := strings.Cut(s, ".")
		_, err := strconv.ParseInt(before[:1], 0, 64)
		if ok && err == nil {
			s = strings.TrimSpace(after)
		}
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		CookTime:     cook,
		Description:  models.Description{Value: description},
		Keywords:     models.Keywords{Values: keywords.String()},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find("title").Text(),
		Yield:        models.Yield{Value: findYield(root.Find(".servings").Text())},
	}, nil
}
