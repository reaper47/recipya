package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"slices"
	"strings"
)

func scrapeGiallozafferano(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err == nil {
		return rs, nil
	}

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Category.Value = root.Find("a[rel='category tag']").First().Text()

	nodes := root.Find(".post-tags a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})
	rs.Keywords.Values = strings.Join(extensions.Unique(keywords), ",")

	node := root.Find(".entry-content").First()
	if len(node.Nodes) > 0 {
		for c := node.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
					if c2.Data == "li" {
						rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(c2.FirstChild.Data))
					}
				}
				continue
			}

			if c.Data == "p" {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(c.FirstChild.Data)))
			}
		}
	}

	rs.Instructions.Values = slices.DeleteFunc(rs.Instructions.Values, func(s models.HowToStep) bool { return s.Text == "" })

	return rs, nil
}
