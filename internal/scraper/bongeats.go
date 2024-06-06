package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"time"
)

func scrapeBongeats(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseLdJSON(root)
	if err != nil {
		return rs, err
	}

	nodes := root.Find(".recipe-ingredients li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = root.Find(".recipe-process li")
	rs.Instructions.Values = make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	parsed, err := time.Parse("Jan 02, 2006", rs.DatePublished)
	if err == nil {
		rs.DatePublished = parsed.Format(time.DateOnly)
	}

	return rs, nil
}
