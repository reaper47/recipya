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

	rs.Name, _ = root.Find("meta[name='title']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='keywords']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Yield.Value = findYield(root.Find(".serving-amount").Children().Last().Text())

	nodes := root.Find(".ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		if v != "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, v)
		}
	})

	nodes = root.Find(".cooking-instructions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.TrimSpace(s.Text())
		if v != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
		}
	})

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
