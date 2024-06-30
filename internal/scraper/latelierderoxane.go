package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeLatelierderoxane(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")

	name := getNameContent(root, "og:title")
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

	getIngredients(&rs, root.Find(".ingredient"))
	getInstructions(&rs, root.Find(".bloc_texte_simple li"))

	return rs, nil
}
