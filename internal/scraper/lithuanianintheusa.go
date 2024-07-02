package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeLithuanianInTheUSA(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Category.Value = root.Find("a[rel=category]").Last().Text()
	rs.Name = getPropertyContent(root, "og:title")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(regex.Digit.FindString(root.Find("em:contains('{')").First().Text()))

	nodes := root.Find("a[rel='tag']")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	getIngredients(&rs, root.Find("h4:contains('{')").Next().NextUntil("ol"))
	getInstructions(&rs, root.Find("ol").First().Find("li"))

	return rs, nil
}
