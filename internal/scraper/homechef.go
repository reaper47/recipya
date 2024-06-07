package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHomechef(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	yield, _ := root.Find("meta[itemprop='recipeYield']").Attr("content")
	rs.Yield.Value = findYield(yield)
	rs.Image.Value, _ = root.Find("div img").First().Attr("data-srcset")
	rs.Name = root.Find("h1").First().Text()

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		s = strings.TrimSpace(strings.TrimPrefix(s, "Info"))
		if s != "" {
			rs.Ingredients.Values = append(rs.Ingredients.Values, s)
		}
	})

	nodes = root.Find("li[itemprop='itemListElement']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.ReplaceAll(sel.Text(), "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		if s != "" {
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		}
	})

	return rs, nil
}
