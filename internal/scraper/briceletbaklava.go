package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBriceletbaklava(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " - ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	content := root.Find(".ob-section-html")
	description := strings.Trim(content.First().Text(), "\n")
	rs.Description.Value = strings.TrimSpace(description)

	nodes := root.Find(".Post-tags a")
	var xk []string
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s == "" {
			return
		}
		xk = append(xk, s)
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	nodes = content.Last().Find("p")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		if s == " " {
			return
		}
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s))
	})

	nodes = content.Last().Find("ul li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(s)))
	})

	return rs, nil
}
