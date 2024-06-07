package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeClosetcooking(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.DatePublished, _ = root.Find("meta[property='article:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Name, _ = root.Find("meta[property='og:title']").Attr("content")
	rs.PrepTime, _ = root.Find("meta[itemprop='prepTime']").Attr("content")
	rs.CookTime, _ = root.Find("meta[itemprop='cookTime']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("li[itemprop='recipeInstructions']")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	nodes = root.Find(".entry-categories a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, sel.Text())
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:       root.Find("span[itemprop='calories']").Text(),
		Carbohydrates:  root.Find("span[itemprop='carbohydrateContent']").Text(),
		Sugar:          root.Find("span[itemprop='sugarContent']").Text(),
		Protein:        root.Find("span[itemprop='proteinContent']").Text(),
		Fat:            root.Find("span[itemprop='fatContent']").Text(),
		SaturatedFat:   root.Find("span[itemprop='saturatedFatContent']").Text(),
		Cholesterol:    root.Find("span[itemprop='cholesterolContent']").Text(),
		Sodium:         root.Find("span[itemprop='sodiumContent']").Text(),
		Fiber:          root.Find("span[itemprop='fiberContent']").Text(),
		TransFat:       root.Find("span[itemprop='transFatContent']").Text(),
		UnsaturatedFat: root.Find("span[itemprop='unsaturatedFatContent']").Text(),
	}

	rs.Yield.Value = findYield(root.Find("span[itemprop='recipeYield']").Text())

	return rs, nil
}
