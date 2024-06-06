package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeDrinkoteket(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[property='og:title']").Last().Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Category.Value, _ = root.Find("meta[property='article:section']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Last().Attr("content")

	nodes := root.Find("ul.ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	nodes = root.Find("#recipe-utrustning .rbs-img-content")
	rs.Tools.Values = make([]models.Tool, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Tools.Values = append(rs.Tools.Values, models.Tool{Name: strings.TrimSpace(sel.Text())})
	})

	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.PrepTime, _ = root.Find("meta[itemprop='prepTime']").Attr("content")
	rs.CookTime, _ = root.Find("meta[itemprop='cookTime']").Attr("content")

	return rs, nil
}
