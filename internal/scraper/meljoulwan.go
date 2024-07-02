package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMeljoulwan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Yield.Value = findYield(root.Find("p:contains('Serves')").Text())
	getIngredients(&rs, root.Find("h5:contains('Ingredients')").Parent().Find("li"))
	getInstructions(&rs, root.Find("h5:contains('Directions')").Parent().Find("p"))

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
