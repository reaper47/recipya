package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRezeptwelt(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name, _ = root.Find("meta[itemprop='name']").Attr("content")
	rs.Category.Value, _ = root.Find("span[itemprop='recipeCategory']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")
	rs.Image.Value, _ = root.Find("img[itemprop='image']").Attr("src")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.DateModified, _ = root.Find("meta[itemprop='dateModified']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s.Text()))
	})

	rs.PrepTime, _ = root.Find("meta[itemprop='performTime']").Attr("content")
	rs.Cuisine.Value, _ = root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[itemprop='keywords']").Attr("content")

	nodes = root.Find("meta[itemprop='tool']")
	rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		t, _ := s.Attr("content")
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(t))
	})

	nodes = root.Find("ol[itemprop='recipeInstructions'] li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	return rs, nil
}
