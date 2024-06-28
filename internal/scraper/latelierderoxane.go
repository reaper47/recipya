package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeLatelierderoxane(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")

	name, _ := root.Find("meta[name='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " - ")
	if found {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Image.Value, _ = root.Find("meta[name='image']").Attr("content")
	rs.DatePublished, _ = root.Find("time[itemprop='datePublished']").Attr("datetime")

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
	rs.PrepTime = prep

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
	rs.CookTime = cook

	split := strings.Split(root.Find("span.titre:contains('Personnes')").Next().Text(), "/")
	if len(split) > 0 {
		rs.Yield.Value = findYield(split[0])
	}

	nodes := root.Find(".ingredient")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".bloc_texte_simple li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
