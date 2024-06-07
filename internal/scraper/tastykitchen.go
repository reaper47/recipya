package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeTastyKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Category.Value = root.Find("a[rel='category tag']").First().Text()
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.PrepTime, _ = root.Find("time[itemprop='prepTime']").Attr("datetime")
	rs.CookTime, _ = root.Find("time[itemprop='cookTime']").Attr("datetime")

	yieldStr, _ := root.Find("input[name='servings']").Attr("value")
	rs.Yield.Value = findYield(yieldStr)

	nodes := root.Find("span[itemprop='ingredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, v)
	})

	nodes = root.Find("span[itemprop='instructions'] p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
