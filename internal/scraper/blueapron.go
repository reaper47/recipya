package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeBlueapron(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Category.Value, _ = root.Find("meta[itemprop='recipeCategory']").Attr("content")
	rs.Cuisine.Value, _ = root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[itemprop='keywords']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[itemprop='image thumbnailUrl']").Attr("content")

	description := root.Find("p[itemprop='description']").Text()
	rs.Description.Value = strings.TrimPrefix(description, "INGREDIENT IN FOCUS")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.Trim(sel.Text(), "\n")
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.Join(strings.Fields(s), " ")
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("div[itemprop='recipeInstructions'] .step-txt")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(strings.Trim(sel.Text(), "\n"))))
	})

	rs.Description = &models.Description{Value: description}
	rs.Name = strings.Trim(root.Find(".ba-recipe-title__main").Text(), "\n")
	rs.NutritionSchema = &models.NutritionSchema{
		Calories: root.Find("span[itemprop='calories']").Text() + " Cals",
	}
	rs.Yield = &models.Yield{Value: findYield(root.Find("span[itemprop='recipeYield']").Text())}

	return rs, nil
}
