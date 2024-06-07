package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMeljoulwan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Yield.Value = findYield(root.Find("p:contains('Serves')").Text())

	nodes := root.Find("h5:contains('Ingredients')").Parent().Find("li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s))
	})

	nodes = root.Find("h5:contains('Directions')").Parent().Find("p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	var category string
	root.Find("div.post-category a").Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "Blog" && category == "" {
			category = s
		}
	})
	rs.Category.Value = category

	var sb strings.Builder
	root.Find("div.post-tage a").Each(func(_ int, sel *goquery.Selection) {
		sb.WriteString(sel.Text())
		sb.WriteString(",")
	})
	rs.Keywords.Values = strings.TrimSuffix(sb.String(), ",")

	return rs, nil
}
