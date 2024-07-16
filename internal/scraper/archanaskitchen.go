package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeArchanasKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	description := root.Find("span[itemprop=description]").Text()
	description = strings.TrimPrefix(description, "\n")
	description = strings.ReplaceAll(description, "\u00a0", " ")
	rs.Description.Value = strings.TrimSpace(description)

	image := root.Find("img[itemprop=image]").AttrOr("src", "")
	rs.Image.Value = "https://www.archanaskitchen.com" + image

	root.Find("li[itemprop=keywords] a").Each(func(_ int, s *goquery.Selection) {
		rs.Keywords.Values += strings.TrimSpace(s.Text()) + ","
	})
	rs.Keywords.Values = strings.TrimSuffix(rs.Keywords.Values, ",")

	getIngredients(&rs, root.Find("li[itemprop=ingredients]"), []models.Replace{
		{"\n", ""},
		{"\t", ""},
		{" , ", ", "},
	}...)

	getInstructions(&rs, root.Find("li[itemprop=recipeInstructions] p"), []models.Replace{
		{"\u00a0", " "},
		{" .", "."},
	}...)

	rs.PrepTime = root.Find("span[itemprop=prepTime]").AttrOr("content", "")
	rs.CookTime = root.Find("span[itemprop=cookTime]").AttrOr("content", "")
	rs.DatePublished = root.Find("span[itemprop=datePublished]").AttrOr("content", "")
	rs.DateModified = root.Find("span[itemprop=dateModified]").AttrOr("content", "")
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield] p").Text())
	rs.Name = root.Find("h1[itemprop=name]").Text()
	rs.Category = &models.Category{Value: root.Find(".recipeCategory a").Text()}
	rs.Cuisine = &models.Cuisine{Value: root.Find("span[itemprop=recipeCuisine] a").Text()}

	return rs, nil
}
