package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeSallysblog(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getNameContent(root, "description")
	rs.Name = strings.ToLower(root.Find("h1").First().Text())

	prep := root.Find("p:contains('Zubereitungszeit')").Next().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(strings.ToLower(prep), "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = split[i]
		}
	}
	rs.PrepTime = prep

	getIngredients(&rs, root.Find(".recipe-description").Next().Find(".hidden").First().Prev().Find("div.text-lg"))
	getInstructions(&rs, root.Find(".recipe-description div p"))

	return rs, nil
}
