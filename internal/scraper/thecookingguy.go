package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeTheCookingGuy(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	date, err := time.Parse("Jan 02, 2006", rs.DatePublished)
	if err == nil {
		rs.DatePublished = date.Format(time.DateOnly)
	}

	nodes := root.Find("div.ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find("div.directions li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	node := root.Find("div.content_indent").Children().Last()
	if node != nil {
		rs.Yield.Value = findYield(node.Text())
	}

	return rs, nil
}
