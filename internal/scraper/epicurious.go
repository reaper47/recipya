package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeEpicurious(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	parseTime := func(s string) string {
		t := "PT"

		parts := strings.Split(s, " ")
		n := len(parts)

		if n == 2 {
			t += parts[0]
			if strings.HasPrefix(strings.ToLower(parts[1]), "min") {
				t += "M"
			} else {
				t += "H"
			}
		} else if n >= 4 {
			t += parts[0] + "H" + parts[2] + "M"
		}

		return t
	}

	if rs.CookTime != "" && !strings.HasPrefix(rs.CookTime, "PT") {
		rs.CookTime = parseTime(rs.CookTime)
	}

	if rs.PrepTime != "" && !strings.HasPrefix(rs.PrepTime, "PT") {
		rs.PrepTime = parseTime(rs.PrepTime)
	}

	if rs.TotalTime != "" && !strings.HasPrefix(rs.TotalTime, "PT") {
		rs.TotalTime = parseTime(rs.TotalTime)
	}

	return rs, nil
}
