package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeMyPlate(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = root.Find(".mp-recipe-full__description p").Text()
	rs.Name = root.Find(".field--name-title").First().Text()
	rs.Image.Value, _ = root.Find(".field--name-field-recipe-image img").Attr("src")

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

	nodes := root.Find(".ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.Join(strings.Fields(s.Text()), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, v)
	})

	nodes = root.Find(".field--name-field-instructions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      root.Find(".total_calories td").Last().Text(),
		Carbohydrates: root.Find(".carbohydrates td").Last().Text(),
		Cholesterol:   root.Find(".cholesterol td").Last().Text(),
		Fat:           root.Find(".total_fat td").Last().Text(),
		Fiber:         root.Find(".dietary_fiber td").Last().Text(),
		Protein:       root.Find(".protein td").Last().Text(),
		SaturatedFat:  root.Find(".saturated_fat td").Last().Text(),
		Sodium:        root.Find(".sodium td").Last().Text(),
		Sugar:         root.Find(".total_sugars td").Last().Text(),
	}

	return rs, nil
}
