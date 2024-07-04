package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeTheGuccha(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Yield.Value = findYield(root.Find(".wprm-recipe-servings").Text())
	rs.Category = &models.Category{Value: strings.ToLower(root.Find(".wprm-recipe-course").Text())}
	rs.Image = &models.Image{Value: root.Find(".wprm-recipe-image img").First().AttrOr("src", "")}

	getIngredients(&rs, root.Find(".wprm-recipe-ingredient"))
	getInstructions(&rs, root.Find(".wprm-recipe-instruction"))

	prep := root.Find(".wprm-recipe-total_time-minutes").Text()
	if prep != "" {
		rs.PrepTime = "PT" + regex.Digit.FindString(prep) + "M"
	}

	return rs, nil
}
