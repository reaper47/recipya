package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeFarmhousedelivery(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	rs.Category.Value = root.Find("a[rel='category tag']").First().Text()

	content := root.Find(".entry-content")
	content.Find("ul li").Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	node := root.Find("p:contains('Instructions')")
	for {
		if node.Nodes == nil || goquery.NodeName(node) == "footer" {
			break
		}

		node = node.Next()
		s := strings.TrimSpace(node.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	}

	return rs, nil
}
