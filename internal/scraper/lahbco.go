package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeLahbco(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getItempropContent(root, "description")
	rs.Image.Value = getItempropContent(root, "image")
	rs.ThumbnailURL.Value = getItempropContent(root, "thumbnailUrl")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")
	rs.Name = getItempropContent(root, "headline")

	nodes := root.Find("span.blog-item-tag a")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, sel.Text())
	})
	rs.Keywords.Values = strings.Join(keywords, ",")

	root.Find("p:contains('INGREDIENTS')").NextAll().Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		_ = s
		switch goquery.NodeName(sel) {
		case "p":
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
		case "ul":
			sel.Find("li").Each(func(_ int, li *goquery.Selection) {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(li.Text()))
			})
		}
	})

	root.Find("ol li").Each(func(_ int, sel *goquery.Selection) {
		sel.Find("p").Each(func(_ int, p *goquery.Selection) {
			s := strings.TrimSpace(p.Text())
			if !slices.ContainsFunc(rs.Instructions.Values, func(h models.HowToItem) bool { return h.Text == s }) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(p.Text()))
			}
		})
	})
	rs.Yield.Value = findYield(regex.Digit.FindString(root.Find("p:contains('Yield')").Text()))

	return rs, nil
}
