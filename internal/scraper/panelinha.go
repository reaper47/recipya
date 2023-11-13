package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapePanelinha(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	node := root.Find("dd:contains('porções')").Text()
	for _, s := range strings.Split(node, " ") {
		yield, err := strconv.ParseInt(s, 10, 16)
		if err == nil {
			rs.Yield.Value = int16(yield)
		}
	}

	nodes := root.Find(".blockIngredientListingsctn ul li")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = s.Text()
	})

	nodes = root.Find("h5:contains('Modo de preparo')").Next().Find("li")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = strings.TrimSuffix(s.Text(), "\u00a0")
	})

	rs.Ingredients = models.Ingredients{Values: ingredients}
	rs.Instructions = models.Instructions{Values: instructions}
	return rs, nil
}
