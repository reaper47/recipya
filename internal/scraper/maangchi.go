package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMaangchi(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	if len(rs.Ingredients.Values) == 0 {
		nodes := root.Find("h2:contains('Ingredients') + ul li")
		rs.Ingredients.Values = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
		})
	}

	if len(rs.Instructions.Values) == 0 {
		nodes := root.Find("h2:contains('Directions') + ol li")
		rs.Instructions.Values = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := sel.Text()
			before, _, ok := strings.Cut(s, "<img")
			if ok {
				s = before
			}
			rs.Instructions.Values = append(rs.Instructions.Values, s)
		})
	}
	return rs, nil
}
