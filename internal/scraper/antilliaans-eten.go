package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeAntilliaansEten(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.Image.Value = root.Find(".post-image img").AttrOr("src", "")
	rs.Name = strings.TrimSpace(root.Find("h1.post-title").Text())
	rs.PrepTime = root.Find(`time[itemprop="prepTime"]`).AttrOr("datetime", "")
	rs.PrepTime = strings.Replace(rs.PrepTime, "0H", "", 1)

	rs.CookTime = root.Find(`time[itemprop="cookTime"]`).AttrOr("datetime", "")
	rs.CookTime = strings.Replace(rs.CookTime, "0H", "", 1)

	ul := root.Find("h2:contains('Ingrediënten')").First().Next()
	getIngredients(&rs, ul.Children())
	getInstructions(&rs, ul.NextUntil("hr"))

	nodes := root.Find(`meta[property="article:tag"]`)
	xk := make([]string, 0, len(nodes.Nodes))
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, sel.AttrOr("content", ""))
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	return rs, nil
}
