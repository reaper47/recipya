package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeBrianLagerstrom(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[itemprop='thumbnailUrl']").Attr("content")

	name, _ := root.Find("meta[itemprop='name']").Attr("content")
	before, _, ok := strings.Cut(name, "—")
	if ok {
		name = before
	}
	rs.Name = strings.TrimSpace(name)

	datePub, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	before, _, ok = strings.Cut(datePub, "T")
	if ok {
		datePub = before
	}
	rs.DatePublished = datePub

	dateMod, _ := root.Find("meta[itemprop='dateModified']").Attr("content")
	before, _, ok = strings.Cut(dateMod, "T")
	if ok {
		dateMod = before
	}
	rs.DateModified = dateMod

	nodes := root.Find("p:contains('▪')")
	if nodes.Length() == 1 {
		// Ingredients
		parts := strings.Split(nodes.Text(), "▪")
		parts = slices.DeleteFunc(parts, func(s string) bool {
			return s == ""
		})
		rs.Ingredients.Values = make([]string, 0, len(parts))
		for _, s := range parts {
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s))
		}

		// Instructions
		nodes = root.Find("div.sqs-html-content p")
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := strings.TrimSpace(sel.Text())
			if strings.HasPrefix(s, "▪") || strings.HasPrefix(s, "*") || strings.HasPrefix(s, "Makes") ||
				strings.HasPrefix(s, "©") || strings.HasPrefix(s, "Privacy") {
				return
			}

			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		})
	} else {
		// Ingredients
		node := root.Find("ul").First()
		for node.Nodes != nil {
			prev := node.Prev()
			if prev.Nodes[0].Data == "p" {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(prev.Text()))
			}

			node.Children().Each(func(_ int, sel *goquery.Selection) {
				if sel.Nodes[0].Data == "li" {
					rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
				}
			})
			node = node.Next()
		}

		// Instructions
		nodes = root.Find("ol li")
		rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := strings.TrimSpace(sel.Text())
			s = strings.ReplaceAll(s, "\u00a0", " ")
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		})
	}

	node := root.Find("p:contains('portion')").First()
	if node != nil {
		rs.Yield.Value = findYield(node.Text())
	}

	return rs, nil
}
