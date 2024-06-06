package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCostco(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")

	keywordsRaw, _ := root.Find("meta[name='description']").Attr("content")
	var keywords strings.Builder
	for _, s := range strings.Split(keywordsRaw, ",") {
		if s != "" {
			keywords.WriteString(s)
		}
	}
	rs.Keywords.Values = keywords.String()

	h1 := root.Find("h1").Last()
	div := h1.Parent()
	name := h1.Text()
	rs.Name = name

	rs.Image.Value, _ = div.Prev().Find("img").Attr("src")

	nodes := div.Find("ul li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = div.Find("p")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	return rs, nil
}
