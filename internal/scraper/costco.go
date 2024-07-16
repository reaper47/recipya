package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCostco(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")

	keywordsRaw := getNameContent(root, "description")
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

	rs.Image.Value = div.Prev().Find("img").AttrOr("src", "")
	getIngredients(&rs, div.Find("ul li"))
	getInstructions(&rs, div.Find("p"))

	return rs, nil
}
