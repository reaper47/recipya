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

	chInstructions := make(chan []string)
	go func() {
		nodes := root.Find("li[itemprop='recipeInstructions']")

		values := make([]string, nodes.Length())
		defer func() {
			_ = recover()
			chInstructions <- values
		}()

		nodes.Each(func(i int, s *goquery.Selection) {
			values[i] = strings.TrimSpace(s.Find(".recipeInstruction__text").Text())
		})
	}()

	name, _ := root.Find("meta[itemprop='name']").Attr("content")
	image, _ := root.Find("meta[itemprop='image']").Attr("content")
	dateModified, _ := root.Find("meta[itemprop='dateModified']").Attr("content")
	datePublished, _ := root.Find("meta[itemprop='datePublished']").Attr("content")
	keywords, _ := root.Find("meta[itemprop='keywords']").Attr("content")
	category, _ := root.Find("meta[itemprop='recipeCategory']").Attr("content")
	cookingMethod, _ := root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	cuisine, _ := root.Find("meta[itemprop='recipeCuisine']").Attr("content")

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          name,
		Image:         models.Image{Value: image},
		Yield:         models.Yield{Value: <-chYield},
		Keywords:      models.Keywords{Values: keywords},
		Category:      models.Category{Value: category},
		CookingMethod: models.CookingMethod{Value: cookingMethod},
		Cuisine:       models.Cuisine{Value: cuisine},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
		Description:   models.Description{Value: <-chDescription},
	}, err
}
