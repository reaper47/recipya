package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBingingWithBabish(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")
	getIngredients(&rs, root.Find("h3:contains('Ingredients')").First().Parent().Find("ul li"))
	getInstructions(&rs, root.Find("h1:contains('Method')").First().Parent().Find("ol li"))

	name := getPropertyContent(root, "og:title")
	before, _, ok := strings.Cut(name, " — ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	return rs, nil
}
