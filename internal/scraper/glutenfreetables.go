package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeGlutenFreeTables(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = strings.TrimSpace(root.Find(".entry-title[itemprop=name]").First().Text())
	rs.Category.Value = strings.TrimSpace(root.Find("a.qodef-e-category").First().Text())
	rs.Yield.Value = findYield(root.Find("input.qodef-quantity-input").AttrOr("value", "1"))
	getTime(&rs, root.Find(".qodef-recipe-prep-time"), true)

	getIngredients(&rs, root.Find(".qodef-ingredients-items"), []models.Replace{
		{"useFields", ""},
	}...)

	nodes := root.Find(".qodef-m-tags-wrapper a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			xk = append(xk, s)
		}
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	root.Find(".qodef-direction-inner").Each(func(_ int, sel *goquery.Selection) {
		name := sel.Find(".qodef-direction-title").Text()
		sel.Children().NextUntil(".qodef-m-completed").Each(func(_ int, p *goquery.Selection) {
			s := strings.Join(strings.Fields(p.Text()), " ")
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s, &models.HowToItem{
				Name: name,
			}))
		})
	})

	return rs, nil
}
