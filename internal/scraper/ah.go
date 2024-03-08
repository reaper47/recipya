package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeAh(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	before, _, ok := strings.Cut(rs.NutritionSchema.Calories, "energie")
	if ok {
		rs.NutritionSchema.Calories = strings.TrimSpace(before)
	}

	before, _, ok = strings.Cut(rs.NutritionSchema.Carbohydrates, "koolhydraten")
	if ok {
		rs.NutritionSchema.Carbohydrates = strings.TrimSpace(before)
	}

	before, _, ok = strings.Cut(rs.NutritionSchema.Fat, "vet")
	if ok {
		rs.NutritionSchema.Fat = strings.TrimSpace(before)
	}

	before, _, ok = strings.Cut(rs.NutritionSchema.Protein, "eiwit")
	if ok {
		rs.NutritionSchema.Protein = strings.TrimSpace(before)
	}

	before, _, ok = strings.Cut(rs.NutritionSchema.SaturatedFat, "waarvan")
	if ok {
		rs.NutritionSchema.SaturatedFat = strings.TrimSpace(before)
	}

	root.Find("div[class^='recipe-header-time']").Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(strings.ToLower(sel.Text()))

		time := "PT"

		parts := strings.Split(s, " ")
		for i := 0; i < len(parts)-1; i++ {
			num, err := strconv.ParseUint(parts[i], 10, 64)
			if err != nil || num == 0 {
				continue
			}

			label := "M"
			if strings.Contains(parts[i+1], "uur") {
				label = "H"
			}

			time += strconv.FormatUint(num, 10) + label
		}

		if strings.Contains(s, "bereiden") {
			rs.PrepTime = time
		} else if strings.Contains(s, "oventijd") {
			rs.CookTime = time
		}
	})
	return rs, nil
}
