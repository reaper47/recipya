package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"strings"
)

func scrapeBodybuilding(root *goquery.Document) (models.RecipeSchema, error) {
	dateModified, _ := root.Find("meta[property='og:updated_time']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")

	description := root.Find(".BBCMS__content--article-description").Text()

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, found := strings.Cut(name, "|")
	if found {
		name = strings.TrimSpace(before)
	}

	var nutrition models.NutritionSchema
	nodes := root.Find(".bb-recipe__meta-nutrient-label")
	nodes.Each(func(_ int, sel *goquery.Selection) {
		switch sel.Text() {
		case "Calories":
			nutrition.Calories = sel.Prev().Text() + " kcal"
		case "Carbs":
			nutrition.Carbohydrates = sel.Prev().Text()
		case "Protein":
			nutrition.Protein = sel.Prev().Text()
		case "Fat":
			nutrition.Fat = sel.Prev().Text()
		}
	})

	nodes = root.Find(".bb-recipe__ingredient-list-item")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", "")
		s = strings.Join(strings.Fields(s), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".bb-recipe__directions-list-item")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", "")
		instructions = append(instructions, strings.TrimSpace(s))
	})

	node := root.Find(".bb-recipe__directions-timing--prep").Find("time")
	prep, _ := node.Attr("datetime")

	node = root.Find(".bb-recipe__directions-timing--cook").Find("time")
	cook, _ := node.Attr("datetime")

	nodes = root.Find(".bb-recipe__topic")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, sel.Text())
	})
	keywords := strings.Join(extensions.Unique(xk), ",")

	return models.RecipeSchema{
		Category:        models.Category{},
		CookTime:        cook,
		CookingMethod:   models.CookingMethod{},
		Cuisine:         models.Cuisine{},
		DateModified:    dateModified,
		DatePublished:   datePublished,
		Description:     models.Description{Value: strings.TrimSpace(strings.Trim(description, "\n"))},
		Keywords:        models.Keywords{Values: keywords},
		Image:           models.Image{},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            name,
		NutritionSchema: nutrition,
		PrepTime:        prep,
		Tools:           models.Tools{},
		Yield:           models.Yield{Value: findYield(root.Find(".bb-recipe__meta-servings .bb-recipe__meta-value-text").Text())},
		URL:             "",
	}, nil
}
