package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeChetnamakan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Name = root.Find(`h1[itemprop=headline]`).Text()

	node := root.Find("strong:contains('Ingredients')")
	if node != nil {
		for c := node.Nodes[0].NextSibling; c != nil; c = c.NextSibling {
			s := strings.TrimSpace(c.Data)
			if s != "br" {
				rs.Ingredients.Values = append(rs.Ingredients.Values, s)
			}
		}
	}

	node = root.Find("strong:contains('Method')")
	if node != nil {
		for c := node.Nodes[0].NextSibling; c != nil; c = c.NextSibling {
			s := strings.TrimSpace(c.Data)
			if s != "br" {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimPrefix(s, "–")))
			}
		}
	}

	return rs, nil
}
