package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeClosetcooking(root *goquery.Document) (models.RecipeSchema, error) {
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	description, _ := root.Find("meta[property='og:description']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	prep, _ := root.Find("meta[itemprop='prepTime']").Attr("content")
	cook, _ := root.Find("meta[itemprop='cookTime']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		instructions = append(instructions, s)
	})

	nodes = root.Find(".entry-categories a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		xk = append(xk, sel.Text())
	})
	keywords := strings.Join(xk, ",")

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		CookTime:      cook,
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		NutritionSchema: models.NutritionSchema{
			Calories:       root.Find("span[itemprop='calories']").Text(),
			Carbohydrates:  root.Find("span[itemprop='carbohydrateContent']").Text(),
			Sugar:          root.Find("span[itemprop='sugarContent']").Text(),
			Protein:        root.Find("span[itemprop='proteinContent']").Text(),
			Fat:            root.Find("span[itemprop='fatContent']").Text(),
			SaturatedFat:   root.Find("span[itemprop='saturatedFatContent']").Text(),
			Cholesterol:    root.Find("span[itemprop='cholesterolContent']").Text(),
			Sodium:         root.Find("span[itemprop='sodiumContent']").Text(),
			Fiber:          root.Find("span[itemprop='fiberContent']").Text(),
			TransFat:       root.Find("span[itemprop='transFatContent']").Text(),
			UnsaturatedFat: root.Find("span[itemprop='unsaturatedFatContent']").Text(),
		},
		PrepTime: prep,
		Yield:    models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())},
	}, nil
}
