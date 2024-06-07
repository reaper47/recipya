package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrape15gram(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.CookTime, _ = root.Find("meta[itemprop='cookTime']").Attr("content")
	rs.PrepTime, _ = root.Find("meta[itemprop='prepTime']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	nodes = root.Find("span[itemprop='keywords']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		keywords = append(keywords, strings.TrimSpace(sel.Text()))
	})

	rs.Description = &models.Description{Value: strings.TrimSpace(root.Find("p[itemprop='description']").Text())}
	rs.Keywords = &models.Keywords{Values: strings.Join(keywords, ",")}
	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.Yield = &models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())}
	return rs, nil
}
