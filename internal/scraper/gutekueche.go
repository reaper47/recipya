package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"regexp"
	"strings"
)

func scrapeGutekueche(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = root.Find("h1").Text()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.URL = getPropertyContent(root, "og:url")
	rs.Image.Value = getPropertyContent(root, "og:image")

	rs.TotalTime = getGuteKuecheRecipeTime(root, `Gesamtzeit`)
	rs.CookTime = getGuteKuecheRecipeTime(root, `Koch & Ruhezeit`)
	rs.PrepTime = getGuteKuecheRecipeTime(root, `Zubereitungszeit`)

	// Keywords are only distinct by CSS class
	root.Find(".recipe-categories .btn-outline").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			rs.Keywords.Values += `,`
		}
		rs.Keywords.Values += s.Text()
	})

	// Tools are only distinct by CSS class
	root.Find(".recipe-categories span.btn-like-box").Each(func(i int, s *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.HowToItem{Type: "HowToTool", Quantity: 1, Text: s.Text()})
	})

	root.Find(".ingredients-table tr").Each(func(iRow int, row *goquery.Selection) {
		thisRow := ""
		row.Children().Each(func(iData int, data *goquery.Selection) {
			if iData > 0 {
				thisRow += " "
			}
			thisRow += strings.TrimSpace(data.Text())
		})
		rs.Ingredients.Values = append(rs.Ingredients.Values, thisRow)
	})

	root.Find(".rezept-preperation li").Each(func(i int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.HowToItem{Type: "HowToStep", Text: s.Text()})
	})

	root.Find(".nutri-block").Each(func(i int, s *goquery.Selection) {
		switch s.Find("header").First().Text() {
		case "kcal":
			rs.NutritionSchema.Calories = s.Find("div").First().Text()
		case "Fett":
			rs.NutritionSchema.Fat = s.Find("div").First().Text()
		case "Eiweiß":
			rs.NutritionSchema.Protein = s.Find("div").First().Text()
		case "Kohlenhydrate":
			rs.NutritionSchema.Carbohydrates = s.Find("div").First().Text()
		}
	})

	rs.Yield.Value = findYield(root.Find(".portions-group input").First().AttrOr("value", "0"))

	return rs, nil
}

// Extracts the time for a given label from the document
func getGuteKuecheRecipeTime(root *goquery.Document, label string) string {
	recipeTimes := root.Find(".recipe-times").Text()
	mSubstring := regexp.MustCompile(`[0-9]+ min. ` + label)

	if len(mSubstring.FindString(recipeTimes)) > 0 {
		return `PT` + regex.Digit.FindString(mSubstring.FindString(recipeTimes)) + `M`
	}
	return ""
}
