package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCuly(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find("span:contains('personen')").Text())
	getIngredients(&rs, root.Find("div.ingredients li"))
	getInstructions(&rs, root.Find("ol li"))
	getTime(&rs, root.Find("span:contains('Voorbereiding')").Parent(), true)
	getTime(&rs, root.Find("span:contains('Kooktijd')").Parent(), false)

	nodes := root.Find("meta[name='cXenseParse:mhu-article_tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.AttrOr("content", "")
		if s != "" {
			xk = append(xk, s)
		}
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	return rs, nil
}
