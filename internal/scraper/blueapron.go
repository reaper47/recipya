package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBlueapron(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Category.Value = getItempropContent(root, "recipeCategory")
	rs.Cuisine.Value = getItempropContent(root, "recipeCuisine")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Keywords.Values = getItempropContent(root, "keywords")
	rs.Image.Value = getItempropContent(root, "image thumbnailUrl")

	description := root.Find("p[itemprop=description]").Text()
	rs.Description.Value = strings.TrimSpace(strings.TrimPrefix(description, "INGREDIENT IN FOCUS"))

	getIngredients(&rs, root.Find("li[itemprop=recipeIngredient]"), []models.Replace{
		{"\n", " "},
		{"useFields", ""},
	}...)
	getInstructions(&rs, root.Find("div[itemprop=recipeInstructions] .step-txt"))

	rs.Name = strings.Trim(root.Find(".ba-recipe-title__main").Text(), "\n")
	rs.NutritionSchema = &models.NutritionSchema{
		Calories: root.Find("span[itemprop=calories]").Text() + " Cals",
	}
	rs.Yield = &models.Yield{Value: findYield(root.Find("span[itemprop=recipeYield]").Text())}

	return rs, nil
}
