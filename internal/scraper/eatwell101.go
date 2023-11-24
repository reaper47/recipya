package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeEatwell101(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")
	keywords, _ := root.Find("meta[name='keywords']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	datePublished, _ := root.Find("meta[property='article:published_time']").Attr("content")
	dateModified, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("h2:contains('Ingredients')").Next().Find("li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := sel.Text()
		ingredients = append(ingredients, s)
	})

	nodes = root.Find("h2:contains('Directions')")
	var instructions []string
	for {
		nodes = nodes.Next()
		s := nodes.Text()
		instructions = append(instructions, s)
		if goquery.NodeName(nodes) != "p" {
			break
		}
	}

	// The below does not work.
	div := root.Find("#recipecardo p.brandon")

	prep := root.Find("i:contains('Prep Time')").Next().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
			break
		}
	}

	cook := root.Find("i:contains('Prep Time')").Next().Text()
	split = strings.Split(cook, " ")
	isMin = strings.Contains(cook, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			cook = "PT" + split[i] + "M"
			break
		}
	}

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		CookTime:      cook,
		DateModified:  dateModified,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: keywords},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: findYield(div.Find("span:contains('servings')").Text())},
	}, nil
}
