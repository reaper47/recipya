package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapeStreetKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	article := root.Find("article").First()
	rs.Name = article.Find("h1").First().Text()
	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time").Attr("content")
	rs.Image.Value, _ = article.Find("img").First().Attr("src")
	rs.Yield.Value = findYield(article.Find(".c-svgicon--servings").Next().Text())

	nodes := article.Find(".ingredients label")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s.Text()))
	})

	nodes = article.Find(".method-step")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.TrimSpace(strings.TrimSuffix(s.Text(), "\n"))
		if v != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
		}
	})

	return rs, nil
}
