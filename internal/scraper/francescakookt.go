package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeFrancescakookt(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	image, _ := root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("meta[property='article:tag']")
	keywords := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		v, _ := sel.Attr("content")
		keywords = append(keywords, v)
	})

	dateMod, _ := root.Find("meta[property='article:modified_time']").Attr("content")
	datePub, _ := root.Find("meta[property='article:published_time']").Attr("content")

	category := root.Find("h2:contains('Categorieën')").Next().Find("a").First().Text()
	before, _, ok := strings.Cut(category, "recepten")
	if ok {
		category = strings.TrimSpace(before)
	}

	var yield int16
	root.Find("h3:contains('personen')").Each(func(_ int, sel *goquery.Selection) {
		if yield != 0 {
			return
		}
		yield = findYield(sel.Text())
	})

	nodes = root.Find("div.dynamic-entry-content ul").Last().Find("li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	var prep string
	s := root.Find("p:contains('Bereidingstijd')").Text()
	if s != "" {
		parts := strings.Split(s, " ")
		if len(parts) == 3 && strings.Contains(s, "minuten") {
			prep = "PT" + regex.Digit.FindString(s) + "M"
		}
	}

	nodes = root.Find("div.dynamic-entry-content ol").Last().Find("li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, strings.TrimSpace(sel.Text()))
	})

	return models.RecipeSchema{
		Category:      models.Category{Value: category},
		DateModified:  dateMod,
		DatePublished: datePub,
		Description:   models.Description{Value: description},
		Keywords:      models.Keywords{Values: strings.ToLower(strings.Join(keywords, ","))},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: yield},
	}, nil
}
