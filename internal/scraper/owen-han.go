package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeOwenhan(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='datePublished']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='dateModified']").Attr("content")
	rs.Name, _ = root.Find("meta[itemprop='headline']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")
	rs.Category.Value = root.Find(".blog-item-category").First().Text()

	content := root.Find("h4:contains('INGREDIENTS')").Parent()

	nodes := content.Find("ul p")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = content.Find("ol p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
