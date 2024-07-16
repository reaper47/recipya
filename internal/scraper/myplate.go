package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
)

func scrapeMyPlate(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = root.Find(".mp-recipe-full__description p").Text()
	rs.Name = root.Find(".field--name-title").First().Text()
	rs.Image.Value = root.Find(".field--name-field-recipe-image img").AttrOr("src", "")

	yieldStr := root.Find(".mp-recipe-full__detail--yield .mp-recipe-full__detail--data").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "\n", "")
	rs.Yield.Value = findYield(strings.TrimSpace(yieldStr))

	cookTimeText := root.Find(".mp-recipe-full__detail--cook-time .mp-recipe-full__detail--data").Text()
	parts := strings.Split(cookTimeText, " ")
	if len(parts) > 1 {
		letter := "M"
		if strings.HasPrefix(parts[1], "hour") {
			letter = "H"
		}
		rs.CookTime = fmt.Sprintf("PT%s%s", parts[0], letter)
	}

	getIngredients(&rs, root.Find(".ingredients li"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find(".field--name-field-instructions li"))

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      regex.Digit.FindString(root.Find(".total_calories td").Last().Text()),
		Carbohydrates: regex.Digit.FindString(root.Find(".carbohydrates td").Last().Text()),
		Cholesterol:   regex.Digit.FindString(root.Find(".cholesterol td").Last().Text()),
		Fat:           regex.Digit.FindString(root.Find(".total_fat td").Last().Text()),
		Fiber:         regex.Digit.FindString(root.Find(".dietary_fiber td").Last().Text()),
		Protein:       regex.Digit.FindString(root.Find(".protein td").Last().Text()),
		SaturatedFat:  regex.Digit.FindString(root.Find(".saturated_fat td").Last().Text()),
		Sodium:        regex.Digit.FindString(root.Find(".sodium td").Last().Text()),
		Sugar:         regex.Digit.FindString(root.Find(".total_sugars td").Last().Text()),
	}

	return rs, nil
}
