package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCestMaFournee(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = root.Find("time.published").AttrOr("datetime", "")
	rs.Image.Value = root.Find("img").First().AttrOr("src", "")

	nodes := root.Find("main div:contains('RÉALISATION')")
	nodes.NextAll().Each(func(_ int, div *goquery.Selection) {
		s := strings.TrimSpace(div.Text())
		if s == "" {
			return
		}

		style := div.Find("span").First().AttrOr("style", "")
		if strings.Contains(style, "color: #c27ba0") {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
		} else if strings.Contains(div.AttrOr("style", ""), "text-align: justify;") {
			s = strings.Join(strings.Fields(s), " ")
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	return rs, nil
}
