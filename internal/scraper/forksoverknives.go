package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeForksOverKnives(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Description.Value = root.Find(".core-paragraph").Text()

	if rs.Category.Value != "" {
		s := strings.TrimSpace(strings.Replace(strings.ToLower(rs.Category.Value), "recipes", "", 1))
		before, _, found := strings.Cut(s, "&")
		if found {
			rs.Category.Value = strings.TrimSpace(before)
		} else {
			rs.Category.Value = s
		}
	}
	return rs, nil
}
