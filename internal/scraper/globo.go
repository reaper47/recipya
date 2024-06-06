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

		p := root.Find("div[itemprop='description'] .content-text__container")
		s = strings.TrimSpace(p.Text())
	}()

	chIngredients := make(chan []string)
	go func() {
		nodes := root.Find("li[itemprop='recipeIngredient']")

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

		yieldStr, _ := root.Find("span[itemprop='recipeYield']").Attr("content")
		i, _ := strconv.ParseInt(yieldStr, 10, 16)
		yield = int16(i)
	}()

	chInstructions := make(chan []models.HowToStep)
	go func() {
		nodes := root.Find("li[itemprop='recipeInstructions']")

		values := make([]models.HowToStep, nodes.Length())
		defer func() {
			_ = recover()
			chInstructions <- values
		}()

		nodes.Each(func(i int, s *goquery.Selection) {
			values[i] = models.NewHowToStep(strings.TrimSpace(s.Find(".recipeInstruction__text").Text()))
		})
	}()

	rs.Name, _ = root.Find("meta[itemprop='name']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[itemprop='image']").Attr("content")
	rs.DateModified, _ = root.Find("meta[itemprop='dateModified']").Attr("content")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[itemprop='keywords']").Attr("content")
	rs.Category.Value, _ = root.Find("meta[itemprop='recipeCategory']").Attr("content")
	rs.CookingMethod.Value, _ = root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	rs.Cuisine.Value, _ = root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	rs.Yield.Value = <-chYield
	rs.Ingredients.Values = <-chIngredients
	rs.Instructions.Values = <-chInstructions
	rs.Description.Value = <-chDescription

	return rs, err
}
