package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeCdKitchen(root *html.Node) (rs models.RecipeSchema, err error) {
	content := getElement(root, "id", "recipepage")

	chCategory := make(chan string)
	go func() {
		var v string
		defer func() {
			_ = recover()
			chCategory <- v
		}()

		node := getElement(root, "class", "current-page")
		node = getElement(node.Parent.Parent.LastChild.PrevSibling.PrevSibling.PrevSibling.FirstChild.NextSibling, "itemprop", "title")
		v = node.FirstChild.Data
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var v models.Ingredients
		defer func() {
			_ = recover()
			chIngredients <- v
		}()

		node := getElement(content, "itemprop", "recipeIngredient").Parent
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "span" {
				s := strings.ReplaceAll(c.FirstChild.Data, "  ", " ")
				s = strings.TrimSpace(s)
				v.Values = append(v.Values, s)
			}
		}
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var v models.Instructions
		defer func() {
			_ = recover()
			chInstructions <- v
		}()

		node := getElement(content, "itemprop", "recipeInstructions").FirstChild
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				v.Values = append(v.Values, strings.TrimSpace(c.Data))
			}
		}
	}()

	chYield := make(chan int16)
	go func() {
		var i int16
		defer func() {
			_ = recover()
			chYield <- i
		}()

		node := getElement(content, "class", "change-servs-input s18")
		yieldStr := getAttr(node, "value")
		yield, err := strconv.Atoi(yieldStr)
		if err == nil {
			i = int16(yield)
		}
	}()

	chNutrition := make(chan models.NutritionSchema)
	go func() {
		node := getElement(content, "itemprop", "nutrition")
		chNutrition <- models.NutritionSchema{
			Calories:       getElement(node, "itemprop", "calories").FirstChild.Data,
			Carbohydrates:  <-getElementData(node, "class", "carbohydrateContent"),
			Sugar:          <-getElementData(node, "class", "sugarContent"),
			Protein:        <-getElementData(node, "class", "proteinContent"),
			Fat:            <-getElementData(node, "class", "fatContent"),
			SaturatedFat:   <-getElementData(node, "class", "saturatedFatContent"),
			Cholesterol:    <-getElementData(node, "class", "cholesterolContent"),
			Sodium:         <-getElementData(node, "class", "sodiumContent"),
			Fiber:          <-getElementData(node, "class", "fiberContent"),
			TransFat:       <-getElementData(node, "class", "transFatContent"),
			UnsaturatedFat: <-getElementData(node, "class", "unsaturatedFatContent"),
		}
	}()

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          models.SchemaType{Value: "Recipe"},
		Category:        models.Category{Value: <-chCategory},
		Name:            <-getItemPropData(content, "name"),
		Description:     models.Description{Value: <-getItemPropData(content, "description")},
		Yield:           models.Yield{Value: <-chYield},
		CookTime:        <-getItemPropAttr(content, "cookTime", "content"),
		Ingredients:     <-chIngredients,
		Instructions:    <-chInstructions,
		NutritionSchema: <-chNutrition,
	}, err
}
