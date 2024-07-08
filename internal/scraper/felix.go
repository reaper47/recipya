package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeFelix(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")

	var xk []string
	root.Find(".cat-links a").Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.ToLower(strings.Join(xk, ","))

	content := root.Find(".entry-content")
	yieldStr := content.Find("h2").Next().Text()
	if yieldStr == "" {
		rs.Yield.Value = findYield(regex.Digit.FindString(content.Find(".wp-block-columns").First().Prev().Text()))
	} else {
		rs.Yield.Value = findYield(yieldStr)
	}

	getIngredients(&rs, content.Find("[style^='color:#008000']"))
	if len(rs.Ingredients.Values) == 0 {
		content.Find(".rp").Each(func(_ int, sel *goquery.Selection) {
			h, _ := sel.Html()
			lines := strings.Split(strings.ReplaceAll(h, "•", ""), "<br/>")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				l = strings.TrimSuffix(strings.TrimPrefix(l, "<em>"), "</em>")
				rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(l))
			}
		})
	}

	root.Find("h2").Next().NextUntil("p:contains('Inspiration: Nagi')").Find("span").Each(func(_ int, sel *goquery.Selection) {
		if strings.Contains(sel.AttrOr("style", ""), "#ffffff") {
			s := strings.Join(strings.Fields(sel.Nodes[0].NextSibling.Data), " ")
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})
	if len(rs.Instructions.Values) == 0 {
		content.Find("div").FilterFunction(func(_ int, sel *goquery.Selection) bool {
			return sel.Has("p.has-drop-cap") != nil
		}).Each(func(_ int, sel *goquery.Selection) {
			s := strings.TrimSpace(sel.Find("p.has-drop-cap").Parent().Text())
			if s == "" {
				return
			}

			s = strings.Join(strings.Fields(s), " ")
			if len(s) > 1 {
				s = s[1:]
			}

			if !slices.ContainsFunc(rs.Instructions.Values, func(item models.HowToItem) bool {
				return item.Text == s
			}) {

				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
			}
		})
	}

	return rs, nil
}
