package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHeatherChristo(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.PrepTime, _ = root.Find("time[itemprop='prepTime']").Attr("datetime")
	rs.CookTime, _ = root.Find("time[itemprop='cookTime']").Attr("datetime")
	rs.Name = strings.TrimSpace(root.Find("div[itemprop='name']").First().Text())
	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").First().Text())

	nodes := root.Find(".ERSIngredients ul").First().Find("li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".ERSInstructions ol").First().Find("li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	return rs, nil
}
