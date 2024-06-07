package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeCooktalk(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished, _ = root.Find("time.entry-date").Attr("datetime")

	nodes := root.Find("a[rel='category']")
	xc := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xc = append(xc, sel.Text())
	})

	rs.Category.Value = xc[0]
	rs.Image.Value, _ = root.Find("img[itemprop='image']").Attr("src")

	description := root.Find("div[itemprop='description']").Text()
	rs.Description.Value = strings.TrimSpace(strings.Trim(description, "\n"))

	nodes = root.Find("li[itemprop='ingredients']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		s = strings.TrimSpace(s)
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("p[itemprop='recipeInstructions']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		s = strings.TrimSpace(s)
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	rs.Name = root.Find(".page-title").Text()

	return rs, nil
}
