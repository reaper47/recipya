package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

func scrapeZeit(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.DateModified = getNameContent(root, "last-modified")
	rs.DatePublished = getNameContent(root, "date")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")

	meta := root.Find(".recipe-list-meta")
	node := meta.Find("title:contains('Portionen')")
	if node.Parent() != nil && node.Parent().Parent() != nil {
		rs.Yield.Value = findYield(node.Parent().Parent().Text())
	}

	getIngredients(&rs, root.Find(".recipe-list-collection__list-item"), []models.Replace{{"useFields", ""}}...)
	getInstructions(&rs, root.Find(".paragraph.article__item"), []models.Replace{{"useFields", ""}}...)

	return rs, nil
}
