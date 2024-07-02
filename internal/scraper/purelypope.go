package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/reaper47/recipya/internal/models"
)

func scrapePurelyPope(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.DatePublished = getPropertyContent(root, "article:published_time")
	image, _ := root.Find("img[itemprop=image]").Attr("src")
	split := strings.Split(image, "?")
	if len(split) > 0 {
		image = split[0]
	}
	rs.Image.Value = image

	rs.Name = root.Find("h2[itemprop=name]").Text()
	rs.Yield.Value = findYield(root.Find("span[itemprop=recipeYield]").Text())

	prepTime, _ := root.Find("time[itemprop=prepTime]").Attr("datetime")
	rs.PrepTime = strings.ReplaceAll(prepTime, " ", "")

	cookTime, _ := root.Find("time[itemprop=cookTime]").Attr("datetime")
	rs.CookTime = strings.ReplaceAll(cookTime, " ", "")

	getIngredients(&rs, root.Find("span[itemprop=recipeIngredient]").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.TrimSpace(s.Text()) != ""
	}))

	getInstructions(&rs, root.Find("div[itemprop=recipeInstructions] li"))

	return rs, nil
}
