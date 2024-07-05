package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeTheHeartySoul(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")

	var (
		isIngredients  = false
		isInstructions = false
	)

	root.Find("#article-content ul").Each(func(_ int, sel *goquery.Selection) {
		prev := sel.Prev()
		prevLower := strings.ToLower(prev.Text())

		if prev.AttrOr("class", "") == "wp-block-heading" {

			if strings.Contains(prevLower, "directions") {
				isIngredients = false
				isInstructions = true
			} else if strings.Contains(prevLower, "ingredients") {
				isIngredients = true
				isInstructions = false
			}
		}

		if isIngredients {
			if !strings.Contains(prevLower, "ingredients") {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(prev.Text()))
			}

			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
			})
		} else if isInstructions {
			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(li.Text()))
			})
		}
	})

	return rs, nil
}
