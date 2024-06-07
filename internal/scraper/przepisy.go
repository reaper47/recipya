package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapePrzepisy(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	rs.Image.Value, _ = root.Find(".recipe-img img").Attr("src")

	yield := findYield(root.Find(".person-count").Text())

	nodes := root.Find(".ingredients-list-content-container")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		name := s.Find(".ingredient-name").Text()
		quantity := s.Find(".quantity").Text()
		ing := fmt.Sprintf("%s %s", quantity, name)
		ingredients[i] = strings.Join(strings.Fields(ing), " ")
	})

	rs.Yield = &models.Yield{Value: yield}
	rs.Ingredients = &models.Ingredients{Values: ingredients}
	return rs, nil
}
