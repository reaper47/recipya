package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeLidl(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = root.Find("title").Text()
	rs.Description.Value = getNameContent(root, "description")
	getTime(&rs, root.Find("span:contains('Bereiding')").Parent(), true)
	rs.Yield.Value = findYield(root.Find("[data-rid='num-servings']").Text())
	rs.Image.Value = root.Find("picture").First().Find("source").Last().AttrOr("srcset", "")

	nodes := root.Find("h4:contains('Tags')").Parent().Find("a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	nodes = root.Find("*[data-rid=ingredient]")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Parent().Text())
		s = strings.ReplaceAll(s, ",", ".")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	getInstructions(&rs, root.Find("*[data-rid=cooking-step] p"))

	return rs, nil
}
