package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeMoulinex(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	node := root.Find("span[aria-labelledby='timePreparationTime']").First()
	if node != nil {
		parts := strings.Split(strings.TrimSpace(node.Text()), " ")
		if len(parts) == 2 {
			h, _, _ := strings.Cut(parts[0], "h")
			m, _, _ := strings.Cut(parts[1], "min")
			rs.PrepTime = "PT" + h + "H" + m + "M"
		} else {
			m, _, _ := strings.Cut(parts[0], "min")
			rs.PrepTime = "PT" + m + "M"
		}
	}

	node = root.Find("span[aria-labelledby='timerCookingTime']").First()
	if node != nil {
		parts := strings.Split(strings.TrimSpace(node.Text()), " ")
		if len(parts) == 2 {
			h, _, _ := strings.Cut(parts[0], "h")
			m, _, _ := strings.Cut(parts[1], "min")
			rs.CookTime = "PT" + h + "H" + m + "M"
		} else {
			m, _, _ := strings.Cut(parts[0], "min")
			rs.CookTime = "PT" + m + "M"
		}
	}

	replace := map[string]string{
		"&#339;": "œ",
	}
	for k, v := range replace {
		rs.Name = strings.ReplaceAll(rs.Name, k, v)
	}

	return rs, nil
}
