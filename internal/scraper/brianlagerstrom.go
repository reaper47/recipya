package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeBrianLagerstrom(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[itemprop='description']").Attr("content")
	image, _ := root.Find("meta[itemprop='thumbnailUrl']").Attr("content")

	name, _ := root.Find("meta[itemprop='name']").Attr("content")
	before, _, ok := strings.Cut(name, "—")
	if ok {
		name = before
	}

	datePublished, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	before, _, ok = strings.Cut(datePublished, "T")
	if ok {
		datePublished = before
	}

	dateModified, _ := root.Find("meta[itemprop='dateModified']").Attr("content")
	before, _, ok = strings.Cut(dateModified, "T")
	if ok {
		dateModified = before
	}

	var (
		ingredients  []string
		instructions []string
	)

	nodes := root.Find("p:contains('▪')")
	if nodes.Length() == 1 {
		// Ingredients
		parts := strings.Split(nodes.Text(), "▪")
		parts = slices.DeleteFunc(parts, func(s string) bool {
			return s == ""
		})
		ingredients = make([]string, 0, len(parts))
		for _, s := range parts {
			ingredients = append(ingredients, strings.TrimSpace(s))
		}

		// Instructions
		nodes = root.Find("div.sqs-html-content p")
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := strings.TrimSpace(sel.Text())
			if strings.HasPrefix(s, "▪") || strings.HasPrefix(s, "*") || strings.HasPrefix(s, "Makes") ||
				strings.HasPrefix(s, "©") || strings.HasPrefix(s, "Privacy") {
				return
			}

			instructions = append(instructions, s)
		})
	} else {
		// Ingredients
		node := root.Find("ul").First()
		for node.Nodes != nil {
			prev := node.Prev()
			if prev.Nodes[0].Data == "p" {
				ingredients = append(ingredients, strings.TrimSpace(prev.Text()))
			}

			node.Children().Each(func(_ int, sel *goquery.Selection) {
				if sel.Nodes[0].Data == "li" {
					ingredients = append(ingredients, sel.Text())
				}
			})
			node = node.Next()
		}

		// Instructions
		nodes = root.Find("ol li")
		instructions = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := strings.TrimSpace(sel.Text())
			s = strings.ReplaceAll(s, "\u00a0", " ")
			instructions = append(instructions, s)
		})
	}

	var yield int16
	node := root.Find("p:contains('portion')").First()
	if node != nil {
		yield = findYield(node.Text())
	}

	return models.RecipeSchema{
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          strings.TrimSpace(name),
		Yield:         models.Yield{Value: yield},
	}, nil
}
