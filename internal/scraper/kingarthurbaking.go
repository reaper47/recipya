package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
	"time"
)

func scrapeKingArthurBaking(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if rs.DatePublished != "" {
		before, _, _ := strings.Cut(rs.DatePublished, "at")

		t, err := time.Parse("January 2, 2006", strings.TrimSpace(before))
		if err == nil {
			rs.DatePublished = t.Format(time.DateOnly)
		}
	}

	var cat string
	root.Find("nav.breadcrumb").First().Find("a").Each(func(_ int, sel *goquery.Selection) {
		href := sel.AttrOr("href", "")
		if href == "/recipes" {
			return
		}

		s := sel.Text()
		before, _, ok := strings.Cut(s, "&")
		if ok {
			s = before
		}
		cat += strings.TrimSpace(s) + ":"
	})

	rs.Category = &models.Category{
		Value: strings.TrimSpace(strings.ToLower(strings.TrimSuffix(cat, ":"))),
	}

	nodes := root.Find("article.recipe__instructions ol li p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if strings.HasPrefix(s, "By\n        ") {
			return
		}
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
	})

	rs.CookingMethod = &models.CookingMethod{}
	rs.Cuisine = &models.Cuisine{}
	rs.Tools = &models.Tools{}

	return rs, nil
}
