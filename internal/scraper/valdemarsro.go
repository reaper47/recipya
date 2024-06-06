package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeValdemarsro(root *goquery.Document) (models.RecipeSchema, error) {
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	start := root.Find("div[itemprop='description']").Children().First()
	description := start.NextUntil(".post-recipe").FilterFunction(func(_ int, selection *goquery.Selection) bool {
		return selection.Nodes[0].Data == "p"
	}).Text()

	name := root.Find("h1[itemprop='headline']").Text()

	yield := findYield(root.Find(".fa-sort").Parent().Text())

	prepTimeStr := root.Find("span:contains('Tid i alt')").Next().Text()
	parts := strings.Split(prepTimeStr, " ")
	var prepTime string
	switch len(parts) {
	case 2:
		prepTime = "PT" + parts[0] + "M"
	case 4:
		prepTime = "PT" + parts[0] + "H" + parts[2] + "M"
	}

	cookTimeStr := root.Find("span:contains('Arbejdstid')").Next().Text()
	parts = strings.Split(cookTimeStr, " ")
	var cookTime string
	switch len(parts) {
	case 2:
		cookTime = "PT" + parts[0] + "M"
	case 4:
		cookTime = "PT" + parts[0] + "H" + parts[2] + "M"
	}

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		ingredients = append(ingredients, s.Text())
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] p")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(s.Text()))
	})

	return models.RecipeSchema{
		Name:          name,
		CookTime:      cookTime,
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: image},
		PrepTime:      prepTime,
		Yield:         &models.Yield{Value: yield},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
	}, nil
}
