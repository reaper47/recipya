package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapePanelinha(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
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
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, s.Text())
	})

	if nodes.Parent() != nil && nodes.Parent().Parent() != nil {
		nodes = root.Find("h5:contains('Modo de preparo')").Parent().Parent().Find("li")
		instructions := make([]models.HowToItem, 0, nodes.Length())
		nodes.Each(func(_ int, s *goquery.Selection) {
			instructions = append(instructions, models.NewHowToStep(strings.TrimSuffix(s.Text(), "\u00a0")))
		})
		rs.Instructions = &models.Instructions{Values: instructions}
	}

	return rs, nil
}
