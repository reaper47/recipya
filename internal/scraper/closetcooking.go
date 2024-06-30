package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
)

func scrapeClosetcooking(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.DateModified = getPropertyContent(root, "article:modified_time")
	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Name = getPropertyContent(root, "og:title")
	rs.PrepTime = getItempropContent(root, "prepTime")
	rs.CookTime = getItempropContent(root, "cookTime")
	rs.Image.Value = getPropertyContent(root, "og:image")

	getIngredients(&rs, root.Find("li[itemprop='recipeIngredient']"))
	getInstructions(&rs, root.Find("li[itemprop='recipeInstructions']"))

	nodes := root.Find(".entry-categories a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, sel.Text())
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:       regex.Digit.FindString(root.Find("span[itemprop='calories']").Text()),
		Carbohydrates:  regex.Digit.FindString(root.Find("span[itemprop='carbohydrateContent']").Text()),
		Sugar:          regex.Digit.FindString(root.Find("span[itemprop='sugarContent']").Text()),
		Protein:        regex.Digit.FindString(root.Find("span[itemprop='proteinContent']").Text()),
		Fat:            regex.Digit.FindString(root.Find("span[itemprop='fatContent']").Text()),
		SaturatedFat:   regex.Digit.FindString(root.Find("span[itemprop='saturatedFatContent']").Text()),
		Cholesterol:    regex.Digit.FindString(root.Find("span[itemprop='cholesterolContent']").Text()),
		Sodium:         regex.Digit.FindString(root.Find("span[itemprop='sodiumContent']").Text()),
		Fiber:          regex.Digit.FindString(root.Find("span[itemprop='fiberContent']").Text()),
		TransFat:       regex.Digit.FindString(root.Find("span[itemprop='transFatContent']").Text()),
		UnsaturatedFat: regex.Digit.FindString(root.Find("span[itemprop='unsaturatedFatContent']").Text()),
	}

	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	return rs, nil
}
