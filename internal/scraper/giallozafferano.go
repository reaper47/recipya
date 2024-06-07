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
	description, _ := root.Find("meta[name='description']").Attr("content")
	rs.Description = &models.Description{Value: description}
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")
	rs.Image = &models.Image{Value: image}
	rs.Category = &models.Category{Value: root.Find("a[rel='category tag']").First().Text()}

	nodes := root.Find(".post-tags a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})
	rs.Keywords = &models.Keywords{Values: strings.Join(extensions.Unique(keywords), ",")}

	var (
		ingredients  []string
		instructions []models.HowToItem
	)

	node := root.Find(".entry-content").First()
	if len(node.Nodes) > 0 {
		for c := node.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
					if c2.Data == "li" {
						ingredients = append(ingredients, strings.TrimSpace(c2.FirstChild.Data))
					}
				}
				continue
			}

			if c.Data == "p" {
				instructions = append(instructions, models.NewHowToStep(strings.TrimSpace(c.FirstChild.Data)))
			}
		}
	}

	instructions = slices.DeleteFunc(instructions, func(s models.HowToItem) bool { return s.Text == "" })

	rs.Ingredients = &models.Ingredients{Values: ingredients}
	rs.Instructions = &models.Instructions{Values: instructions}
	rs.Yield = &models.Yield{Value: 1}
	return rs, nil
}
