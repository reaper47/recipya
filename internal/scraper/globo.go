package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGlobo(root *goquery.Document) (rs models.RecipeSchema, err error) {
	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		p := root.Find("div[itemprop=description] .content-text__container")
		s = strings.TrimSpace(p.Text())
	}()

	chIngredients := make(chan []string)
	go func() {
		nodes := root.Find("li[itemprop=recipeIngredient]")

		values := make([]string, nodes.Length())
		defer func() {
			_ = recover()
			chIngredients <- values
		}()

		nodes.Each(func(i int, s *goquery.Selection) {
			values[i] = s.Text()
		})
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		yieldStr := root.Find("span[itemprop=recipeYield]").AttrOr("content", "")
		i, _ := strconv.ParseInt(yieldStr, 10, 16)
		yield = int16(i)
	}()

	chInstructions := make(chan []models.HowToItem)
	go func() {
		nodes := root.Find("li[itemprop=recipeInstructions]")

		values := make([]models.HowToItem, nodes.Length())
		defer func() {
			_ = recover()
			chInstructions <- values
		}()

		nodes.Each(func(i int, s *goquery.Selection) {
			values[i] = models.NewHowToStep(strings.TrimSpace(s.Find(".recipeInstruction__text").Text()))
		})
	}()

	rs.Name = getItempropContent(root, "name")
	rs.Image = &models.Image{Value: getItempropContent(root, "image")}
	rs.DateModified = getItempropContent(root, "dateModified")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.Keywords = &models.Keywords{Values: getItempropContent(root, "keywords")}
	rs.Category = &models.Category{Value: getItempropContent(root, "recipeCategory")}
	rs.CookingMethod = &models.CookingMethod{Value: getItempropContent(root, "recipeCuisine")}
	rs.Cuisine = &models.Cuisine{Value: getItempropContent(root, "recipeCuisine")}
	rs.Yield = &models.Yield{Value: <-chYield}
	rs.Ingredients = &models.Ingredients{Values: <-chIngredients}
	rs.Instructions = &models.Instructions{Values: <-chInstructions}
	rs.Description = &models.Description{Value: <-chDescription}

	return rs, err
}
