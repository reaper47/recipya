package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeKuchynalidla(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Yield.Value = findYield(root.Find(".recipe-detail._servings").Text())

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, "|")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")

	var prep string
	matches := slices.DeleteFunc(regex.Time.FindStringSubmatch(root.Find(".recipe-detail._time").First().Text()), func(s string) bool {
		return s == ""
	})
	switch len(matches) {
	case 3:
		prep = "PT" + regex.Digit.FindString(matches[1]) + "H" + regex.Digit.FindString(matches[2]) + "M"
	case 2:
		v := regex.Digit.FindString(matches[1])
		if strings.Contains(matches[1], "h") {
			prep = "PT" + v + "H"
		} else {
			prep = "PT" + v + "M"
		}
	}
	rs.PrepTime = prep

	nodes := root.Find("div.ing ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	node := root.Find("h2:contains('Postup')")
	if len(node.Nodes) > 0 {
		for c := node.Nodes[0]; c != nil; c = c.NextSibling {
			if c.Data == "h2" || c.Data == "p" {
				s := strings.TrimSpace(c.FirstChild.Data)
				if s == "" {
					continue
				}
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
			} else if c.Data == "ul" {
				if c.PrevSibling == nil && c.PrevSibling.PrevSibling.FirstChild == nil {
					continue
				}

				if strings.Contains(strings.ToUpper(c.PrevSibling.PrevSibling.FirstChild.Data), "POTREBUJEME") {
					goquery.NewDocumentFromNode(c).Children().Each(func(_ int, sel *goquery.Selection) {
						rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(strings.TrimSpace(sel.Text())))
					})
				}
			}
		}
	}

	return rs, nil
}
