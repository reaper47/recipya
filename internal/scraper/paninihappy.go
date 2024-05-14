package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePaniniHappy(root *goquery.Document) (models.RecipeSchema, error) {
	content := root.Find(".entry-content")
	image, _ := content.Find("img").First().Attr("src")

	recipe := content.Find(".hrecipe")

	var description string
	content.Children().NextUntil(".hrecipe").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			description += "\n\n"
		}
		description += s.Text()
	})
	description = strings.TrimSuffix(description, "\n\n\n")

	var prepTime string
	prepTimeStr := recipe.Find(".preptime").Text()
	parts := strings.Split(prepTimeStr, " ")
	if len(parts) > 1 {
		letter := "M"
		if strings.HasPrefix(parts[1], "hour") {
			letter = "H"
		}
		prepTime = fmt.Sprintf("PT%s%s", parts[0], letter)
	}

	var cookTime string
	cookeTimeStr := recipe.Find(".cooktime").Text()
	parts = strings.Split(cookeTimeStr, " ")
	if len(parts) > 1 {
		letter := "M"
		if strings.HasPrefix(parts[1], "hour") {
			letter = "H"
		}
		cookTime = fmt.Sprintf("PT%s%s", parts[0], letter)
	}

	yield := findYield(recipe.Find(".yield").Text())

	nodes := recipe.Find(".ingredient")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = s.Text()
	})

	nodes = recipe.Find(".instruction")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = s.Text()
	})

	return models.RecipeSchema{
		Name:         recipe.Find("h2").Last().Text(),
		Description:  models.Description{Value: description},
		Image:        models.Image{Value: image},
		PrepTime:     prepTime,
		CookTime:     cookTime,
		Yield:        models.Yield{Value: yield},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
	}, nil
}
