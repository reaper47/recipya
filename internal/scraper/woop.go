package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeWoop(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[name='title']").Attr("content")
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
			nutrition.Calories = val
		case "Protein":
			nutrition.Protein = val
		case "Carbohydrate":
			nutrition.Carbohydrates = val
		case "Fat":
			nutrition.Fat = val
		}
	})
	rs.NutritionSchema = &nutrition

	return rs, nil
}
