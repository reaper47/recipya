package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePaniniHappy(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	content := root.Find(".entry-content")
	rs.Image.Value, _ = content.Find("img").First().Attr("src")

	recipe := content.Find(".hrecipe")
	rs.Name = recipe.Find("h2").Last().Text()

	var description string
	content.Children().NextUntil(".hrecipe").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			description += "\n\n"
		}
		description += s.Text()
	})
	rs.Description.Value = strings.TrimSuffix(description, "\n\n\n")

	prepTimeStr := recipe.Find(".preptime").Text()
	parts := strings.Split(prepTimeStr, " ")
	if len(parts) > 1 {
		letter := "M"
		if strings.HasPrefix(parts[1], "hour") {
			letter = "H"
		}
		rs.PrepTime = fmt.Sprintf("PT%s%s", parts[0], letter)
	}

	cookeTimeStr := recipe.Find(".cooktime").Text()
	parts = strings.Split(cookeTimeStr, " ")
	if len(parts) > 1 {
		letter := "M"
		if strings.HasPrefix(parts[1], "hour") {
			letter = "H"
		}
		rs.CookTime = fmt.Sprintf("PT%s%s", parts[0], letter)
	}

	rs.Yield.Value = findYield(recipe.Find(".yield").Text())

	nodes := recipe.Find(".ingredient")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, s.Text())
	})

	nodes = recipe.Find(".instruction")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
