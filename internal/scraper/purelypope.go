package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePurelyPope(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	image, _ := root.Find("img[itemprop='image']").Attr("src")
	split := strings.Split(image, "?")
	if len(split) > 0 {
		image = split[0]
	}
	rs.Image.Value = image

	rs.Name = root.Find("h2[itemprop='name']").Text()
	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	prepTime, _ := root.Find("time[itemprop='prepTime']").Attr("datetime")
	rs.PrepTime = strings.ReplaceAll(prepTime, " ", "")

	cookTime, _ := root.Find("time[itemprop='cookTime']").Attr("datetime")
	rs.CookTime = strings.ReplaceAll(cookTime, " ", "")

	nodes := root.Find("span[itemprop='recipeIngredient']").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.TrimSpace(s.Text()) != ""
	})
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(s.Text()))
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s.Text()))
	})

	return rs, nil
}
