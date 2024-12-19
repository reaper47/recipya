package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMindMegette(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Ingredients.Values = make([]string, 0)
	root.Find("div.ingredients").Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Find(".ingredients-group").Text())

		sel.Find(".ingredients-meta").Each(func(_ int, meta *goquery.Selection) {
			var ingredient []string
			meta.Children().Each(func(_ int, meta *goquery.Selection) {
				ingredient = append(ingredient, meta.Text())
			})
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.Join(strings.Fields(strings.Join(ingredient, " ")), " "))
		})
	})

	nodes := root.Find("div.right-side").First().Find("ol").First().Find("li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	rs.Name = root.Find(".recipe-title").Text()
	rs.Yield.Value = findYield(root.Find(".recipe-details div:contains('dag')").Text())
	return rs, nil
}
