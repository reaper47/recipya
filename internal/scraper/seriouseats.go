package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeSeriousEats(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err == nil {
		return rs, nil
	}

	rs = models.NewRecipeSchema()
	rs.Name = strings.TrimSpace(strings.TrimSuffix(getPropertyContent(root, "og:title"), "Recipe"))
	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = getPropertyContent(root, "og:image")

	getTime(&rs, root.Find("div.cook-time"), false)
	getTime(&rs, root.Find("div.prep-time"), true)

	nodes := root.Find("ul.tag-nav__list li a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	getIngredients(&rs, root.Find("ul.structured-ingredients__list li"))
	getInstructions(&rs, root.Find("#section--instructions_1-0 ol li > p"))

	yieldStr := regex.Digit.FindString(root.Find("div.recipe-serving").Text())
	yield, err := strconv.ParseInt(yieldStr, 10, 16)
	if err != nil {
		yield = 0 // or handle the error as appropriate
	}
	rs.Yield.Value = int16(yield)

	root.Find("#toc-special-equipment").Next().Next().Children().Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(sel.Text()))
	})

	nutrition := root.Find("div.nutrition-label")
	if nutrition.Length() > 0 {
		extractNutrition := func(nutrition string) string {
			return regex.Digit.FindString(strings.ReplaceAll(extensions.ConvertToString(nutrition), ",", "."))
		}

		rs.NutritionSchema.Calories = nutrition.Find("td:contains('Calories')").Next().Text()
		rs.NutritionSchema.Fat = extractNutrition(nutrition.Find("td:contains('Total Fat')").Text())
		rs.NutritionSchema.SaturatedFat = extractNutrition(nutrition.Find("td:contains('Saturated Fat')").Text())
		rs.NutritionSchema.Cholesterol = extractNutrition(nutrition.Find("td:contains('Cholesterol')").Text())
		rs.NutritionSchema.Sodium = extractNutrition(nutrition.Find("td:contains('Sodium')").Text())
		rs.NutritionSchema.Carbohydrates = extractNutrition(nutrition.Find("td:contains('Total Carbohydrate')").Text())
		rs.NutritionSchema.Fiber = extractNutrition(nutrition.Find("td:contains('Fiber')").Text())
		rs.NutritionSchema.Sugar = extractNutrition(nutrition.Find("td:contains('Sugar')").Text())
		rs.NutritionSchema.Protein = extractNutrition(nutrition.Find("td:contains('Protein')").Text())
	}

	return rs, nil
}
