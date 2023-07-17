package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeMyPlate(root *goquery.Document) (models.RecipeSchema, error) {
	description := root.Find(".mp-recipe-full__description p").Text()

	name := root.Find(".field--name-title").First().Text()
	image, _ := root.Find(".field--name-field-recipe-image img").Attr("src")

	yieldStr := root.Find(".mp-recipe-full__detail--yield .mp-recipe-full__detail--data").Text()
	yieldStr = strings.ReplaceAll(yieldStr, "\n", "")
	yield := findYield(strings.TrimSpace(yieldStr))

	cookTimeText := root.Find(".mp-recipe-full__detail--cook-time .mp-recipe-full__detail--data").Text()
	parts := strings.Split(cookTimeText, " ")
	letter := "M"
	if strings.HasPrefix(parts[1], "hour") {
		letter = "H"
	}
	cookTime := fmt.Sprintf("PT%s%s", parts[0], letter)

	nodes := root.Find(".ingredients li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.Join(strings.Fields(s.Text()), " ")
		ingredients[i] = v
	})

	nodes = root.Find(".field--name-field-instructions li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		AtContext:    atContext,
		AtType:       models.SchemaType{Value: "Recipe"},
		CookTime:     cookTime,
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         name,
		NutritionSchema: models.NutritionSchema{
			Calories:      root.Find(".total_calories td").Last().Text(),
			Carbohydrates: root.Find(".carbohydrates td").Last().Text(),
			Cholesterol:   root.Find(".cholesterol td").Last().Text(),
			Fat:           root.Find(".total_fat td").Last().Text(),
			Fiber:         root.Find(".dietary_fiber td").Last().Text(),
			Protein:       root.Find(".protein td").Last().Text(),
			SaturatedFat:  root.Find(".saturated_fat td").Last().Text(),
			Sodium:        root.Find(".sodium td").Last().Text(),
			Sugar:         root.Find(".total_sugars td").Last().Text(),
		},
		PrepTime: "",
		Yield:    models.Yield{Value: yield},
	}, nil
}
