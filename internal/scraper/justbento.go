package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeJustbento(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='og:updated_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='og:published_time']").Attr("content")

	category := root.Find("nav.breadcrumb").Find("a:contains('Recipe collection:')").Text()
	_, after, found := strings.Cut(category, ":")
	if found {
		category = strings.TrimSpace(after)
	} else {
		category = ""
	}

	nodes := root.Find(".field-name-body li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".field-name-body ul").Last()
	for {
		nodes = nodes.Next()
		if nodes.Nodes == nil {
			break
		}

		if goquery.NodeName(nodes) != "p" {
			continue
		}

		s := strings.TrimSpace(nodes.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	}

	rs.Image.Value, _ = root.Find(".field-name-body img").First().Attr("src")
	rs.Cuisine.Value = "Japanese"
	rs.Category.Value = category
	rs.Yield.Value = findYield(root.Find("*:contains('portions')").Text())

	return rs, nil
}
