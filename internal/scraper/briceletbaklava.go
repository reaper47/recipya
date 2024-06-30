package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBriceletbaklava(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name := getPropertyContent(root, "og:title")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Image.Value = getPropertyContent(root, "og:image")

	content := root.Find(".ob-section-html")
	description := strings.Trim(content.First().Text(), "\n")
	rs.Description.Value = strings.TrimSpace(description)

	nodes := root.Find(".Post-tags a")
	var xk []string
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s != "" {
			xk = append(xk, s)
		}
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	getIngredients(&rs, content.Last().Find("p"))
	getInstructions(&rs, content.Last().Find("ul li"))

	return rs, nil
}
