package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeYumelise(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value, _ = root.Find(".wprm-recipe-image").First().Find("img").First().Attr("src")
	rs.Name = root.Find(".wprm-recipe-name").First().Text()
	rs.Category.Value = root.Find(".wprm-recipe-course-container").First().Children().Last().Text()
	rs.Cuisine.Value = root.Find(".wprm-recipe-cuisine-container").First().Children().Last().Text()
	rs.Description.Value = root.Find(".wprm-recipe-summary").First().Text()

	prep := root.Find(".wprm-recipe-prep_time-minutes").First().Text()
	if prep != "" {
		prep, _, _ = strings.Cut(prep, " ")
		prep = "PT" + prep + "M"
	}
	rs.PrepTime = prep

	cook := root.Find(".wprm-recipe-cook_time-minutes").First().Text()
	if cook != "" {
		cook, _, _ = strings.Cut(cook, " ")
		cook = "PT" + cook + "M"
	}
	rs.CookTime = cook

	yield, _ := root.Find(".wprm-recipe-servings").First().Attr("data-original-servings")
	rs.Yield.Value = findYield(yield)

	var keywords []string
	root.Find("a[rel='tag']").Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, sel.Text())
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	getIngredients(&rs, root.Find(".wprm-recipe-ingredient"))
	getInstructions(&rs, root.Find(".wprm-recipe-instruction-text"), []models.Replace{{"\u00a0", " "}}...)

	return rs, nil
}
