package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeMaangchi(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	node := root.Find("h2:contains('Ingredients')")
	if len(node.Nodes) > 0 {
		clear(rs.Ingredients.Values)

		for c := node.Nodes[0].NextSibling; c != nil; c = c.NextSibling {
			if c.Data == "h2" && strings.ToLower(c.FirstChild.Data) == "directions" {
				break
			}

			switch c.Data {
			case "h3":
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(c.FirstChild.Data))
			case "ul":
				goquery.NewDocumentFromNode(c).Find("li").Each(func(_ int, sel *goquery.Selection) {
					rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
				})
			}
		}
	}

	node = root.Find("h2:contains('Directions')")
	if node != nil {
		clear(rs.Instructions.Values)

		for c := node.Nodes[0].NextSibling; c != nil; c = c.NextSibling {
			if c.Data == "h2" && strings.ToLower(c.FirstChild.Data) == "directions" {
				break
			}

			switch c.Data {
			case "h3":
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(c.FirstChild.Data)))
			case "p":
				s := strings.TrimSpace(c.FirstChild.Data)
				if s != "a" {
					rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
				}
			case "ol":
				goquery.NewDocumentFromNode(c).Find("li").Each(func(_ int, sel *goquery.Selection) {
					s := sel.Text()
					before, _, ok := strings.Cut(s, "<img")
					if ok {
						s = before
					}
					rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(s)))
				})
			}
		}
	}

	rs.Ingredients.Values = slices.DeleteFunc(rs.Ingredients.Values, func(s string) bool { return s == "" })
	rs.Instructions.Values = slices.DeleteFunc(rs.Instructions.Values, func(s models.HowToStep) bool { return s.Text == "" })
	return rs, nil
}
