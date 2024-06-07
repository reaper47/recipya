package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRosannapansino(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.DatePublished, _ = root.Find(".article__date time").First().Attr("datetime")
	rs.Image.Value, _ = root.Find("div.article__body img").First().Attr("src")
	rs.Yield.Value = findYield(root.Find("p:contains('Makes')").Text())

	description, _ := root.Find("meta[property='og:description']").Attr("content")
	before, _, found := strings.Cut(strings.TrimSpace(description), "\n")
	if found {
		description = strings.TrimSpace(before)
	}
	rs.Description.Value = description

	content := root.Find(".recipe-content")

	nodes := content.Find("ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = content.Find("ol li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
