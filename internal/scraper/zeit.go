package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeZeit(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.DateModified, _ = root.Find("meta[name='last-modified']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[name='date']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='keywords']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")

	meta := root.Find(".recipe-list-meta")

	node := meta.Find("title:contains('Portionen')")
	if node.Parent() != nil && node.Parent().Parent() != nil {
		rs.Yield.Value = findYield(node.Parent().Parent().Text())
	}

	nodes := root.Find(".recipe-list-collection__list-item")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Join(strings.Fields(sel.Text()), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".paragraph.article__item")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.Join(strings.Fields(s), " ")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	return rs, nil
}
