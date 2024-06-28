package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeHomebrewAnswers(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err != nil {
		rs = models.NewRecipeSchema()
	}

	if rs.Category.Value == "uncategorized" {
		cat := strings.ToLower(root.Find(".blog-categories").First().Text())
		before, _, ok := strings.Cut(cat, "recipe")
		if ok {
			rs.Category.Value = strings.TrimSpace(before)
		} else {
			rs.Category.Value = strings.TrimSpace(cat)
		}
	} else {
		rs.Category.Value = strings.ToLower(rs.Category.Value)
	}

	if rs.DatePublished == "" {
		rs.DatePublished = getPropertyContent(root, "article:published_time")
	}

	if rs.DateModified == "" {
		rs.DateModified = getPropertyContent(root, "article:modified_time")
	}

	if rs.Image.Value == "" {
		rs.Image.Value = getPropertyContent(root, "og:image")
	}

	rs.Description.Value = getPropertyContent(root, "og:description")

	if len(rs.Ingredients.Values) == 0 {
		nodes := root.Find("h2:contains('Ingredients')").Next().Find("li")
		rs.Ingredients.Values = make([]string, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
		})
	} else {
		rs.Ingredients.Values = slices.DeleteFunc(rs.Ingredients.Values, func(s string) bool { return s == "" })

	}

	if rs.CookingMethod == nil {
		rs.CookingMethod = &models.CookingMethod{}
	}

	if rs.Keywords == nil {
		rs.Keywords = &models.Keywords{}
	}

	if rs.NutritionSchema == nil {
		rs.NutritionSchema = &models.NutritionSchema{}
	}

	if rs.ThumbnailURL == nil {
		rs.ThumbnailURL = &models.ThumbnailURL{}
	}

	if rs.Tools == nil || len(rs.Tools.Values) == 0 {
		rs.Tools = &models.Tools{}

		nodes := root.Find("h3:contains('Equipment')")
		if nodes.Length() == 0 {
			nodes = root.Find("h2:contains('gallon / ')")
		}

		nodes = nodes.Next().Find("li")
		rs.Tools.Values = make([]models.HowToItem, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(sel.Text()))
		})
	}

	if rs.Name == "" {
		rs.Name = strings.TrimSpace(root.Find(".post-title").First().Text())
	}

	if len(rs.Instructions.Values) == 0 {
		nodes := root.Find("h3:contains('Method')")
		if nodes.Length() == 0 {
			nodes = root.Find("h2:contains('Method')")
			if nodes.Length() == 0 {
				nodes = root.Find("h2:contains('Ingredients')").Next().NextAll()
			}
		}

		nodes = nodes.NextAll()
		rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
		nodes.Each(func(_ int, sel *goquery.Selection) {
			s := sel.Text()
			dotIndex := strings.Index(s, ".")
			if dotIndex != -1 && dotIndex < 4 {
				_, s, _ = strings.Cut(s, ".")
			}
			rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
		})
	}

	return rs, nil
}
