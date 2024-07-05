package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"time"
)

func scrapeBongeats(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	getIngredients(&rs, root.Find(".recipe-ingredients li"))
	getInstructions(&rs, root.Find(".recipe-process li"))

	parsed, err := time.Parse("Jan 02, 2006", rs.DatePublished)
	if err == nil {
		rs.DatePublished = parsed.Format(time.DateOnly)
	}

	return rs, nil
}
