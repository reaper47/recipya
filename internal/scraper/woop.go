package scraper

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
)

func scrapeWoop(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getNameContent(root, "title")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find(".serving-amount").Children().Last().Text())

	getIngredients(&rs, root.Find(".ingredients li"))
	getInstructions(&rs, root.Find(".cooking-instructions li"))

	var nutrition models.NutritionSchema
	root.Find(".nutritional-info li").Each(func(_ int, s *goquery.Selection) {
		parts := strings.Split(s.Text(), ":")
		val := strings.TrimSpace(strings.Join(parts[1:], " "))
		switch parts[0] {
		case "Energy":
			if match := regexp.MustCompile(`(\d+)Kcal`).FindStringSubmatch(val); len(match) > 1 {
				nutrition.Calories = match[1]
			} else {
				nutrition.Calories = ""
			}
		case "Protein":
			nutrition.Protein = regex.Digit.FindString(val)
		case "Carbohydrate":
			nutrition.Carbohydrates = regex.Digit.FindString(val)
		case "Fat":
			nutrition.Fat = regex.Digit.FindString(val)
		}
	})
	rs.NutritionSchema = &nutrition

	return rs, nil
}
