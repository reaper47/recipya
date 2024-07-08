package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeDonnaHay(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Description.Value = getNameContent(root, "description")
	rs.Name = strings.TrimSpace(strings.TrimSuffix(getPropertyContent(root, "og:title"), "| Donna Hay"))
	rs.Image.Value = getPropertyContent(root, "og:image")

	getIngredients(&rs, root.Find(".ingredients ul li"))
	getInstructions(&rs, root.Find(".method").First().Find("ol li"))

	cat := root.Find("link[rel=canonical]").AttrOr("href", "")
	if cat != "" {
		before, _, _ := strings.Cut(strings.TrimPrefix(cat, "https://www.donnahay.com.au/recipes/"), "/")
		rs.Category.Value = strings.ReplaceAll(before, "-", " ")
	}
	return rs, nil
}
