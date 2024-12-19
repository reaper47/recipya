package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/extensions"
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

	getTime(&rs, root.Find("span:contains('Zubereitungszeit')"), true)
	getIngredients(&rs, root.Find(".shop-studio-recipes-recipe-detail-tabs-description-ingredients__content__ingredient-list__ingredient"))
	rs.Ingredients.Values = extensions.Unique(rs.Ingredients.Values)
	getInstructions(&rs, root.Find(".shop-studio-recipes-recipe-detail-tabs-description-preparations__content__preparation__text"))

	return rs, nil
}
