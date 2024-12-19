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
	getTime(&rs, root.Find("p:contains('Prep Time')"), true)
	getTime(&rs, root.Find("p:contains('Cook Time')"), false)

	nodes := root.Find("meta[property='article:tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.AttrOr("content", "")
		if s != "" {
			xk = append(xk, s)
		}
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	cat := strings.ToLower(strings.TrimSpace(root.Find("p:contains('Course:')").Text()))
	if cat != "" {
		rs.Category.Value = strings.TrimSpace(strings.TrimPrefix(cat, "course:"))
	}

	root.Find("h3:contains('Ingredients')").NextUntil("h3").Each(func(_ int, sel *goquery.Selection) {
		switch goquery.NodeName(sel) {
		case "p":
			rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
		case "ul":
			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, li.Text())
			})
		}
	})

	root.Find("h3:contains('Equipment')").Next().Children().Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(sel.Text()))
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
