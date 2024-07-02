package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeTheFoodFlamingo(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DateModified = getPropertyContent(root, "og:updated_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.Name = root.Find("#recipe").Next().Find("h2.wp-block-heading").First().Text()

	nodes := root.Find("meta[property='article:tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s, _ := sel.Attr("content")
		xk = append(xk, s)
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	prep := strings.TrimSpace(root.Find("p:contains('Prep Time')").Text())
	if prep != "" {
		unit := "H"
		if strings.Contains(prep, "min") {
			unit = "M"
		}
		rs.PrepTime = "PT" + regex.Digit.FindString(prep) + unit
	}

	cook := strings.TrimSpace(root.Find("p:contains('Cook Time')").Text())
	if cook != "" {
		unit := "H"
		if strings.Contains(cook, "min") {
			unit = "M"
		}
		rs.CookTime = "PT" + regex.Digit.FindString(cook) + unit
	}

	cat := strings.ToLower(strings.TrimSpace(root.Find("p:contains('Course:')").Text()))
	if cat != "" {
		rs.Category.Value = strings.TrimSpace(strings.TrimPrefix(cat, "course:"))
	}

	root.Find("ul").FilterFunction(func(_ int, sel *goquery.Selection) bool {
		_, exists := sel.Attr("class")
		return !exists && goquery.NodeName(sel.Prev()) != "h2"
	}).Each(func(_ int, ul *goquery.Selection) {
		isEquipment := strings.Contains(ul.Prev().Text(), "Equipment")
		ul.Children().Each(func(_ int, li *goquery.Selection) {
			s := strings.TrimSpace(li.Text())
			if isEquipment {
				rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(s))
			} else {
				rs.Ingredients.Values = append(rs.Ingredients.Values, s)
			}
		})
	})

	root.Find("ol").Each(func(_ int, ol *goquery.Selection) {
		name := strings.TrimSpace(ol.Prev().Text())
		ol.Children().Each(func(_ int, li *goquery.Selection) {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(li.Text(), &models.HowToItem{Name: name}))
		})
	})

	if root.Find("h3:contains('Nutrition')").Text() != "" {
		rs.NutritionSchema = &models.NutritionSchema{
			Calories:      regex.Digit.FindString(root.Find("p:contains('Calories:')").Text()),
			Carbohydrates: regex.Digit.FindString(root.Find("p:contains('Carbs:')").Text()),
			Fat:           regex.Digit.FindString(root.Find("p:contains('Total Fat:')").Text()),
			Fiber:         regex.Digit.FindString(root.Find("p:contains('Fibers:')").Text()),
			Protein:       regex.Digit.FindString(root.Find("p:contains('Protein:')").Text()),
			Sodium:        regex.Digit.FindString(root.Find("p:contains('Sodium:')").Text()),
			Sugar:         regex.Digit.FindString(root.Find("p:contains('Sugars:')").Text()),
		}
	}

	return rs, nil
}
