package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeBbcgoodfood(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		return rs, err
	}

	rs.NutritionSchema.Servings = strconv.Itoa(int(rs.Yield.Value))

	before, _, found := strings.Cut(rs.NutritionSchema.Carbohydrates, "carbohydrates")
	if found {
		rs.NutritionSchema.Carbohydrates = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Sugar, "sugar")
	if found {
		rs.NutritionSchema.Sugar = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Protein, "protein")
	if found {
		rs.NutritionSchema.Protein = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Fat, "fat")
	if found {
		rs.NutritionSchema.Fat = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.SaturatedFat, "saturated fat")
	if found {
		rs.NutritionSchema.SaturatedFat = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.UnsaturatedFat, "unsaturated fat")
	if found {
		rs.NutritionSchema.UnsaturatedFat = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Sodium, "of sodium")
	if found {
		rs.NutritionSchema.Sodium = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Fiber, "fiber")
	if found {
		rs.NutritionSchema.Fiber = strings.TrimSpace(before)
	}

	before, _, found = strings.Cut(rs.NutritionSchema.Cholesterol, "cholesterol")
	if found {
		rs.NutritionSchema.Cholesterol = strings.TrimSpace(before)
	}

	return rs, nil
}
