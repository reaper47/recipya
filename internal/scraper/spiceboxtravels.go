package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeSpiceBoxTravels(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(root.Find("p:contains('Serves:')").Text())

	nodes := root.Find("[rel=tag]")
	if nodes.Length() > 0 {
		xk := make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			xk = append(xk, strings.TrimSpace(sel.Text()))
		})

		rs.Category.Value = xk[0]
		if len(xk) > 1 {
			rs.Keywords.Values = strings.Join(xk[1:], ",")
		}
	}

	var isIngredients = true
	root.Find(".entry-content p").Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s == "" || (isIngredients && len(s) > 60) {
			return
		}

		dotIndex := strings.Index(s, ".")
		if dotIndex != -1 && dotIndex < 4 {
			isIngredients = false
		}

		if !isIngredients {
			_, after, _ := strings.Cut(s, ".")
			if after != "" {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(after))
			}
		} else if isIngredients {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
		}
	})

	return rs, nil
}
