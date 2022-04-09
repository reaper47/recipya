package scraper

import (
	"strconv"
	"strings"
	"sync"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeAfghanKitchen(root *html.Node) (rs models.RecipeSchema, err error) {
	content := getElement(root, "id", "content")

	var wg sync.WaitGroup
	wg.Add(2)
	var info *html.Node
	go func() {
		defer wg.Done()
		info = getElement(content, "class", "recipe-info clearfix")
	}()
	var about *html.Node
	go func() {
		defer wg.Done()
		about = getElement(content, "itemprop", "about")
	}()
	wg.Wait()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(info, "class", "servings")
		node = getElement(node, "class", "value")
		i, err := strconv.Atoi(node.FirstChild.Data)
		if err == nil {
			yield = int16(i)
		}
	}()

	chPrepTime := make(chan string)
	go func() {
		var prep string
		defer func() {
			_ = recover()
			chPrepTime <- prep
		}()

		node := getElement(info, "class", "prep-time")
		node = getElement(node, "class", "value")

		time := node.FirstChild.Data
		if strings.Contains(time, "m") {
			prep = "PT" + strings.TrimSuffix(time, "m") + "M"
		} else if strings.Contains(time, "h") {
			time = strings.TrimSuffix(time, " h")
			parts := strings.Split(time, ":")
			prep = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}()

	chCookTime := make(chan string)
	go func() {
		var cook string
		defer func() {
			_ = recover()
			chCookTime <- cook
		}()

		node := getElement(info, "class", "cook-time")
		node = getElement(node, "class", "value")

		time := node.FirstChild.Data
		if strings.Contains(time, "m") {
			cook = "PT" + strings.TrimSuffix(time, "m") + "M"
		} else if strings.Contains(time, "h") {
			time = strings.TrimSuffix(time, " h")
			parts := strings.Split(time, ":")
			cook = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}()

	description := about.FirstChild.NextSibling.FirstChild.Data
	description = strings.ReplaceAll(description, "\n", "")
	description = strings.ReplaceAll(description, "\u00a0", " ")

	chIngredients := make(chan []string)
	chInstructions := make(chan []string)
	go func() {
		var vals []string
		ingNode := getElement(about, "class", "blue")
		for c := ingNode.NextSibling.FirstChild; c != nil; c = c.NextSibling {
			ing := strings.ReplaceAll(c.FirstChild.Data, "  ", " ")
			vals = append(vals, ing)
		}
		chIngredients <- vals

		vals = []string{}
		node := ingNode.NextSibling.NextSibling
		for {
			if node == nil {
				break
			}

			if node.Data == "p" {
				ins := strings.TrimSpace(node.FirstChild.Data)
				vals = append(vals, ins)
			}

			node = node.NextSibling
		}
		chInstructions <- vals
	}()

	chCategory := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chCategory <- s
		}()

		node := getElement(about, "class", "type")
		s = node.FirstChild.NextSibling.FirstChild.Data
	}()

	return models.RecipeSchema{
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-getItemPropData(content, "name"),
		DatePublished: <-getItemPropAttr(content, "datePublished", "content"),
		Image:         models.Image{Value: <-getItemPropAttr(content, "image", "content")},
		Yield:         models.Yield{Value: <-chYield},
		PrepTime:      <-chPrepTime,
		CookTime:      <-chCookTime,
		Description:   models.Description{Value: strings.TrimSpace(description)},
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
		Category:      models.Category{Value: <-chCategory},
	}, nil
}
