package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeWoop(root *goquery.Document) (rs models.RecipeSchema, err error) {
	name, _ := root.Find("meta[name='title']").Attr("content")
	keywords, _ := root.Find("meta[name='keywords']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	yield := findYield(root.Find(".serving-amount").Children().Last().Text())

	nodes := root.Find(".ingredients li")
	var ingredients []string
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		if v != "" {
			ingredients = append(ingredients, v)
		}
	})

	instructions := make([]string, 0)
	root.Find(".cooking-instructions li").Each(func(i int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		if v != "" {
			instructions = append(instructions, v)
		}
	})

	var nutrition models.NutritionSchema
	root.Find(".nutritional-info li").Each(func(i int, s *goquery.Selection) {
		parts := strings.Split(s.Text(), ":")
		val := strings.TrimSpace(strings.Join(parts[1:], " "))
		switch parts[0] {
		case "Energy":
			nutrition.Calories = val
		case "Protein":
			nutrition.Protein = val
		case "Carbohydrate":
			nutrition.Carbohydrates = val
		case "Fat":
			nutrition.Fat = val
		}
	})

	return models.RecipeSchema{
		AtContext:       atContext,
		AtType:          models.SchemaType{Value: "Recipe"},
		Name:            name,
		Description:     models.Description{Value: description},
		Image:           models.Image{Value: image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Keywords:        models.Keywords{Values: keywords},
		NutritionSchema: nutrition,
		Yield:           models.Yield{Value: yield},
	}, nil
}
