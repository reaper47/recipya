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

	node := root.Find("h2:contains('Ingredients:')")
	node.NextUntil("h3:contains('Directions:')").Each(func(_ int, sel *goquery.Selection) {
		switch goquery.NodeName(sel) {
		case "p":
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
		case "ul":
			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
			})
		}
	})

	node = root.Find("h3:contains('Directions:')")
	node.NextAll().Each(func(_ int, sel *goquery.Selection) {
		if goquery.NodeName(sel) == "ul" {
			sel.Children().Each(func(_ int, li *goquery.Selection) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(li.Text()))
			})
		}
	})

	return rs, nil
}
