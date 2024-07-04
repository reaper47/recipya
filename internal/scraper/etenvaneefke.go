package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapEetenvaneefke(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Name = getPropertyContent(root, "og:title")

	nodes := root.Find(".entry-content ul")
	getIngredients(&rs, nodes.Find("li"))
	rs.Ingredients.Values = slices.DeleteFunc(rs.Ingredients.Values, func(s string) bool {
		return strings.HasPrefix(s, "Klik om")
	})

	getInstructions(&rs, nodes.NextAll().Filter("p"))

	return rs, nil
}
