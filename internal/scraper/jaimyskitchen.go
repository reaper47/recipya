package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeJaimysKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs.Keywords.Values = strings.Join(slices.DeleteFunc(strings.Split(rs.Keywords.Values, ","), func(s string) bool {
		return strings.TrimSpace(s) == ""
	}), ",")

	rs.PrepTime = rs.TotalTime
	return rs, nil
}
