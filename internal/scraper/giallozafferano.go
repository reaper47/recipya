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

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	description, _ := root.Find("meta[name='description']").Attr("content")
	datePub, _ := root.Find("meta[property='article:published_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find(".post-tags a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(strings.ToLower(sel.Text())))
	})

	var (
		ingredients  []string
		instructions []string
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
				instructions = append(instructions, strings.TrimSpace(c.FirstChild.Data))
			}
		}
	}

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Category:      models.Category{Value: root.Find("a[rel='category tag']").First().Text()},
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: strings.Join(extensions.Unique(keywords), ",")},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: slices.DeleteFunc(instructions, func(s string) bool { return s == "" })},
		Name:          name,
		Yield:         models.Yield{Value: 1},
	}, nil
}
