package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeCdKitchen(root *goquery.Document) (rs models.RecipeSchema, err error) {
	category := root.Find(".prev-page").Last().Text()
	category = strings.Join(strings.Fields(category), " ")

	content := root.Find("#recipepage")

	nodes := content.Find("span[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "  ", " ")
		ingredients[i] = strings.TrimSpace(v)
	})

	node := content.Find("div[itemprop='recipeInstructions'] p")
	node.Find("br").Each(func(_ int, s *goquery.Selection) {
		s.ReplaceWithHtml("$$$")
	})
	lines := strings.Split(node.Text(), "$$$$$$")
	instructions := make([]string, len(lines))
	for i, line := range lines {
		instructions[i] = strings.TrimSpace(line)
	}

	yieldStr, _ := content.Find(".change-servs-input").Attr("value")
	yield, _ := strconv.ParseInt(yieldStr, 10, 16)

	node = content.Find("span[itemprop='nutrition']")
	nutrition := models.NutritionSchema{
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

	cookTime, _ := content.Find("meta[itemprop='cookTime']").Attr("content")

	return models.RecipeSchema{
		AtContext:       atContext,
		AtType:          models.SchemaType{Value: "Recipe"},
		Category:        models.Category{Value: category},
		Name:            content.Find("h1[itemprop='name']").Text(),
		Description:     models.Description{Value: content.Find("p[itemprop='description']").Text()},
		Yield:           models.Yield{Value: int16(yield)},
		CookTime:        cookTime,
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		NutritionSchema: nutrition,
	}, err
}
