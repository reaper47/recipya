package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strconv"
	"strings"
)

func scrapeQuitoque(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = root.Find("#sylius-product-name").Text()

	productTagsNode := root.Find("#product-tags")

	var keywords []string
	productTagsNode.Find("span").Each(func(_ int, s *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(s.Text()))
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	rs.Description.Value = strings.TrimSpace(productTagsNode.Next().Text())

	ingredientsNode := root.Find("#ingredients ul")
	rs.Yield.Value = findYield(ingredientsNode.Prev().Children().Last().Text())

	nodes := root.Find("#ingredients .ingredient-list li")
	nodes = nodes.AddSelection(root.Find("ul.kitchen-list li"))
	getIngredients(&rs, nodes, []models.Replace{{"useFields", ""}}...)

	getInstructions(&rs, root.Find("#steps-default > div").Children(), models.Replace{
		Old: "\n                                    ",
		New: "\n",
	})
	rs.Instructions.Values = slices.DeleteFunc(rs.Instructions.Values, func(item models.HowToItem) bool {
		return len(item.Text) < 10
	})

	root.Find("#equipment ul li").Each(func(_ int, s *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(s.Text()))
	})

	extractTime := func(value string) string {
		parts := strings.Split(value, " ")
		var s string
		if len(parts) == 2 {
			s = "PT"
			if strings.HasPrefix(parts[1], "min") {
				s += parts[0] + "M"
			} else {
				s += parts[0] + "H"
			}
		} else if len(parts) == 4 {
			s = "PT" + parts[0] + "H" + parts[1] + "M"
		}
		return s
	}

	root.Find(".recipe-infos-short .item-info").Each(func(_ int, sel *goquery.Selection) {
		value := strings.TrimSpace(sel.Find("p").First().Text())

		text := strings.TrimSpace(sel.Find(".regular").Text())
		if strings.Contains(text, "cuisine") {
			rs.PrepTime = extractTime(value)
		} else if strings.HasPrefix(text, "Total") {
			totalTimeStr := extractTime(value)
			var (
				totalTime int
				err       error
			)
			if strings.Contains(totalTimeStr, "H") {
				totalTimeStr = strings.TrimPrefix(totalTimeStr, "PT")
				parts := strings.Split(totalTimeStr, "H")
				minutes, err := strconv.Atoi(parts[1])
				if err != nil {
					return
				}

				hours, err := strconv.Atoi(parts[0])
				if err != nil {
					return
				}

				totalTime = minutes + (hours * 60)
			} else {
				totalTime, err = strconv.Atoi(regex.Digit.FindString(totalTimeStr))
				if err != nil {
					return
				}
			}

			var prepTime int
			prepTimeStr := strings.TrimPrefix(rs.PrepTime, "PT")
			if strings.Contains(prepTimeStr, "H") {
				parts := strings.Split(prepTimeStr, "H")
				minutes, err := strconv.Atoi(parts[1])
				if err != nil {
					return
				}

				hours, err := strconv.Atoi(parts[0])
				if err != nil {
					return
				}

				prepTime = minutes + (hours * 60)
			} else {
				prepTime, err = strconv.Atoi(strings.TrimSuffix(prepTimeStr, "M"))
				if err != nil {
					return
				}
			}

			rs.CookTime = "PT" + strconv.Itoa(totalTime-prepTime) + "M"
		}
	})

	extractNutrition := func(nutrition string) string {
		return regex.Digit.FindString(strings.ReplaceAll(extensions.ConvertToString(nutrition), ",", "."))
	}

	root.Find("#quantity ul li").Each(func(_ int, sel *goquery.Selection) {
		name := strings.ToLower(sel.Children().First().Text())
		v := extractNutrition(sel.Children().Last().Text())

		switch name {
		case "énergie (kcal)":
			rs.NutritionSchema.Calories = v
		case "matières grasses":
			rs.NutritionSchema.Fat = v
		case "dont acides gras saturés":
			rs.NutritionSchema.SaturatedFat = v
		case "glucides":
			rs.NutritionSchema.Carbohydrates = v
		case "dont sucre":
			rs.NutritionSchema.Sugar = v
		case "fibres":
			rs.NutritionSchema.Fiber = v
		case "protéines":
			rs.NutritionSchema.Protein = v
		case "sel":
			rs.NutritionSchema.Sodium = v
		}
	})

	return rs, nil
}
