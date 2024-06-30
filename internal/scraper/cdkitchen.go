package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeCdKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	category := root.Find(".prev-page").Last().Text()
	rs.Category.Value = strings.Join(strings.Fields(category), " ")

	content := root.Find("#recipepage")

	getIngredients(&rs, content.Find("span[itemprop='recipeIngredient']"), []models.Replace{
		{"  ", " "},
	}...)

	node := content.Find("div[itemprop='recipeInstructions'] p")
	node.Find("br").Each(func(_ int, s *goquery.Selection) {
		s.ReplaceWithHtml("$$$")
	})
	lines := strings.Split(node.Text(), "$$$$$$")
	rs.Instructions.Values = make([]models.HowToItem, 0, len(lines))
	for _, line := range lines {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(line))
	}

	yieldStr, _ := content.Find(".change-servs-input").Attr("value")
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)
	rs.Yield.Value = int16(yield)

	node = content.Find("span[itemprop='nutrition']")
	rs.NutritionSchema = &models.NutritionSchema{
		Calories:       node.Find("span[itemprop='calories']").Text(),
		Carbohydrates:  node.Find(".carbohydrateContent").Text(),
		Sugar:          node.Find(".sugarContent").Text(),
		Protein:        node.Find(".proteinContent").Text(),
		Fat:            node.Find(".fatContent").Text(),
		SaturatedFat:   node.Find(".saturatedFatContent").Text(),
		Cholesterol:    node.Find(".cholesterolContent").Text(),
		Sodium:         node.Find(".sodiumContent").Text(),
		Fiber:          node.Find(".fiberContent").Text(),
		TransFat:       node.Find(".transFatContent").Text(),
		UnsaturatedFat: node.Find(".unsaturatedFatContent").Text(),
	}

	rs.CookTime, _ = content.Find("meta[itemprop='cookTime']").Attr("content")
	rs.Description = &models.Description{Value: content.Find("p[itemprop='description']").Text()}
	rs.Name = content.Find("h1[itemprop='name']").Text()

	return rs, nil
}
