package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeAngielaEats(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getItempropContent(root, "description")
	rs.Image.Value = getItempropContent(root, "image")
	rs.ThumbnailURL.Value = getItempropContent(root, "thumbnailUrl")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")
	rs.Name = getItempropContent(root, "headline")
	getIngredients(&rs, root.Find("ul").Last().Find("li"))

	nodes := root.Find("ol").Last().Find("li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		var s []string
		sel.Find("p").Each(func(_ int, sel2 *goquery.Selection) {
			s = append(s, strings.TrimSpace(sel2.Text()))
		})
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.Join(s, "\n\n")))
	})

	return rs, nil
}
