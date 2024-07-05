package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCooktalk(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = root.Find(".page-title").Text()
	rs.DatePublished, _ = root.Find("time.entry-date").Attr("datetime")

	nodes := root.Find("a[rel='category']")
	xc := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xc = append(xc, sel.Text())
	})

	rs.Category.Value = xc[0]
	rs.Image.Value, _ = root.Find("img[itemprop=image]").Attr("src")

	description := root.Find("div[itemprop=description]").Text()
	rs.Description.Value = strings.TrimSpace(strings.Trim(description, "\n"))

	getIngredients(&rs, root.Find("li[itemprop=ingredients]"))
	getInstructions(&rs, root.Find("p[itemprop=recipeInstructions]"))

	return rs, nil
}
