package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRosannapansino(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.DatePublished = root.Find(".article__date time").First().AttrOr("datetime", "")
	rs.Image.Value = root.Find("div.article__body img").First().AttrOr("src", "")
	rs.Yield.Value = findYield(root.Find("p:contains('Makes')").Text())

	description := getPropertyContent(root, "og:description")
	before, _, found := strings.Cut(strings.TrimSpace(description), "\n")
	if found {
		description = strings.TrimSpace(before)
	}
	rs.Description.Value = description

	content := root.Find(".recipe-content")
	getIngredients(&rs, content.Find("ul li"))
	getInstructions(&rs, content.Find("ol li"))

	return rs, nil
}
