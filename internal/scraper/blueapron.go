package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBlueapron(root *goquery.Document) (models.RecipeSchema, error) {
	category, _ := root.Find("meta[itemprop='recipeCategory']").Attr("content")
	cuisine, _ := root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	datePublished, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	keywords, _ := root.Find("meta[itemprop='keywords']").Attr("content")
	image, _ := root.Find("meta[itemprop='image thumbnailUrl']").Attr("content")

	description := root.Find("p[itemprop='description']").Text()
	description = strings.TrimPrefix(description, "INGREDIENT IN FOCUS")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := strings.Trim(sel.Text(), "\n")
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] .step-txt")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(strings.Trim(sel.Text(), "\n")))
	})

	return models.RecipeSchema{
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: category},
		Cuisine:       models.Cuisine{Value: cuisine},
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          strings.Trim(root.Find(".ba-recipe-title__main").Text(), "\n"),
		NutritionSchema: models.NutritionSchema{
			Calories: root.Find("span[itemprop='calories']").Text() + " Cals",
		},
		Yield: models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())},
	}, nil
}
