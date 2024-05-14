package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeYumelise(root *goquery.Document) (models.RecipeSchema, error) {
	image, _ := root.Find(".wprm-recipe-image").First().Find("img").First().Attr("src")

	prep := root.Find(".wprm-recipe-prep_time-minutes").First().Text()
	if prep != "" {
		prep, _, _ = strings.Cut(prep, " ")
		prep = "PT" + prep + "M"
	}

	cook := root.Find(".wprm-recipe-cook_time-minutes").First().Text()
	if cook != "" {
		cook, _, _ = strings.Cut(cook, " ")
		cook = "PT" + cook + "M"
	}

	yield, _ := root.Find(".wprm-recipe-servings").First().Attr("data-original-servings")

	var keywords []string
	root.Find("a[rel='tag']").Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, sel.Text())
	})

	nodes := root.Find(".wprm-recipe-ingredient")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".wprm-recipe-instruction-text")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(strings.TrimSpace(sel.Text()), "\u00a0", " ")
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		Category:     models.Category{Value: root.Find(".wprm-recipe-course-container").First().Children().Last().Text()},
		CookTime:     cook,
		Cuisine:      models.Cuisine{Value: root.Find(".wprm-recipe-cuisine-container").First().Children().Last().Text()},
		Description:  models.Description{Value: root.Find(".wprm-recipe-summary").First().Text()},
		Keywords:     models.Keywords{Values: strings.Join(keywords, ",")},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         root.Find(".wprm-recipe-name").First().Text(),
		PrepTime:     prep,
		Yield:        models.Yield{Value: findYield(yield)},
	}, nil
}
