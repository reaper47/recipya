package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeHomeChef(root *goquery.Document) (models.RecipeSchema, error) {
	content := root.Find("main")

	yieldStr, _ := content.Find("meta[itemprop='recipeYield']").Attr("content")

	image, _ := content.Find("img").First().Attr("src")

	nodes := content.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = s.Text()
	})

	nodes = content.Find(".meal__steps li[itemprop='description']")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:   "https://schema.org",
		AtType:      models.SchemaType{Value: "Recipe"},
		Image:       models.Image{Value: image},
		Name:        content.Find("span[itemprop='name'] h1").Text(),
		Description: models.Description{Value: content.Find("div[itemprop='description'] p").Text()},
		Yield:       models.Yield{Value: findYield(yieldStr)},
		NutritionSchema: models.NutritionSchema{
			Calories:       content.Find("strong[itemprop='calories']").Text(),
			Carbohydrates:  content.Find("strong[itemprop='carbohydrateContent']").Text(),
			Sugar:          content.Find("strong[itemprop='sugarContent']").Text(),
			Protein:        content.Find("strong[itemprop='proteinContent']").Text(),
			Fat:            content.Find("strong[itemprop='fatContent']").Text(),
			SaturatedFat:   content.Find("strong[itemprop='saturatedFatContent']").Text(),
			Cholesterol:    content.Find("strong[itemprop='cholesterolContent']").Text(),
			Sodium:         content.Find("strong[itemprop='sodiumContent']").Text(),
			Fiber:          content.Find("strong[itemprop='fiberContent']").Text(),
			TransFat:       content.Find("strong[itemprop='transFatContent']").Text(),
			UnsaturatedFat: content.Find("strong[itemprop='unsaturatedFatContent']").Text(),
		},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
	}, nil
}
