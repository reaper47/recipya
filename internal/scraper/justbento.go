package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeJustbento(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getPropertyContent(root, "og:title")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DateModified, _ = root.Find("meta[property='og:updated_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='og:published_time']").Attr("content")

	category := root.Find("nav.breadcrumb").Find("a:contains('Recipe collection:')").Text()
	_, after, found := strings.Cut(category, ":")
	if found {
		category = strings.TrimSpace(after)
	} else {
		category = ""
	}

	getIngredients(&rs, root.Find(".field-name-body li"))

	nodes := root.Find(".field-name-body ul").Last()
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
