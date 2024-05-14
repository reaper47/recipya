package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeLatelierderoxane(root *goquery.Document) (models.RecipeSchema, error) {
	description, _ := root.Find("meta[name='description']").Attr("content")

	name, _ := root.Find("meta[name='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " - ")
	if found {
		name = strings.TrimSpace(before)
	}

	image, _ := root.Find("meta[name='image']").Attr("content")
	datePublished, _ := root.Find("time[itemprop='datePublished']").Attr("datetime")

	prep := root.Find("span:contains('Préparation')").Next().Text()
	if prep != "" {
		split := strings.Split(prep, " ")
		isMin := strings.Contains(prep, "min")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				prep = "PT" + split[i] + "M"
				break
			}
		}
	}

	cook := root.Find("span:contains('Cuisson')").Next().Text()
	if cook != "" {
		split := strings.Split(cook, " ")
		isMin := strings.Contains(cook, "min")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				cook = "PT" + split[i] + "M"
				break
			}
		}
	}

	split := strings.Split(root.Find("span.titre:contains('Personnes')").Next().Text(), "/")
	var yield string
	if len(split) > 0 {
		yield = split[0]
	}

	nodes := root.Find(".ingredient")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".bloc_texte_simple li")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		CookTime:      cook,
		DatePublished: datePublished,
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          name,
		PrepTime:      prep,
		Yield:         models.Yield{Value: findYield(yield)},
	}, nil
}
