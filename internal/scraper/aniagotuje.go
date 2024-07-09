package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeAniagotuje(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Name = getPropertyContent(root, "og:title")
	rs.Image.Value = getPropertyContent(root, "og:image")

	info := root.Find(".recipe-info")

	n := info.Find("strong:contains('przygo')")
	if len(n.Nodes) > 0 && n.Nodes[0].NextSibling != nil {
		prep := strings.TrimSpace(n.Nodes[0].NextSibling.Data)
		if prep != "" {
			rs.PrepTime = "PT" + regex.Digit.FindString(prep)
			if strings.Contains(prep, "min") {
				rs.PrepTime += "M"
			} else {
				rs.PrepTime += "H"
			}
		}
	}

	n = info.Find("strong:contains('gotow')")
	if len(n.Nodes) > 0 && n.Nodes[len(n.Nodes)-1].NextSibling != nil {
		cook := strings.TrimSpace(n.Nodes[len(n.Nodes)-1].NextSibling.Data)
		if cook != "" {
			rs.CookTime = "PT" + regex.Digit.FindString(cook)
			if strings.Contains(cook, "min") {
				rs.CookTime += "M"
			} else {
				rs.CookTime += "H"
			}
		}
	}

	root.Find(".recipe-ing-list").Each(func(_ int, ul *goquery.Selection) {
		header := strings.TrimSpace(ul.Prev().Text())
		if header != "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, header)
		}

		ul.Children().Each(func(_ int, li *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
		})
	})

	root.Find("h3").First().NextUntil("h3").Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	return rs, nil
}
