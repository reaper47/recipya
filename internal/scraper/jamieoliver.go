package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strings"
)

func scrapeJamieOliver(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseWebsite(root)
	if err == nil {
		return rs, nil
	}

	rs = models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Keywords.Values = getNameContent(root, "keywords")
	rs.Image.Value = getPropertyContent(root, "og:image")
	rs.Yield.Value = findYield(strings.TrimSpace(root.Find(".recipe-detail.serves").Text()))
	rs.Name = root.Find("h1").First().Text()

	getIngredients(&rs, root.Find("ul.ingred-list").Children(), []models.Replace{
		{"useFields", ""},
	}...)

	root.Find("div.method-p > div").Each(func(i int, sel *goquery.Selection) {
		if sel.HasClass("promo-banner") || sel.HasClass("avocado-desktop-foot") ||
			sel.HasClass("mobile-related-recipes") {
			return
		}

		if sel.HasClass("tip") {
			sel.Find("p").Each(func(_ int, p *goquery.Selection) {
				h, _ := p.Html()
				h = strings.ReplaceAll(h, "<br/>", "\n")
				h = strings.ReplaceAll(h, "  ", "")
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(h))
			})

		} else {
			h, _ := sel.Html()
			parts := strings.Split(h, "<br/>")
			for _, part := range parts {
				v := strings.TrimSpace(part)
				if v != "" {
					rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(v))
				}
			}
		}
	})

	_, after, ok := strings.Cut(root.Find("link[rel=canonical]").AttrOr("href", ""), "https://www.jamieoliver.com/recipes/")
	if ok {
		cat, _, _ := strings.Cut(after, "/")
		rs.Category.Value = strings.TrimSpace(strings.ReplaceAll(strings.TrimSuffix(cat, "recipes"), "-", " "))
	}

	cooksIn := strings.TrimSpace(strings.ToLower(root.Find(".recipe-detail.time").Text()))
	if cooksIn != "" {
		var (
			prep string
			cook string
		)

		if strings.Contains(cooksIn, "prep") {
			_, prep, _ = strings.Cut(cooksIn, "prep")
		}

		if strings.Contains(prep, "cook") {
			_, cook, _ = strings.Cut(prep, "cook")
		}

		if prep != "" {
			rs.PrepTime = "PT" + regex.Digit.FindString(prep)
			if strings.Contains(prep, "min") {
				rs.PrepTime += "M"
			} else {
				rs.PrepTime += "H"
			}
		} else {
			matches := regex.Time.FindAllString(prep, 2)
			if matches != nil {
				matches = slices.DeleteFunc(matches, func(s string) bool { return s == "" })
			}

			if len(matches) == 2 {
				rs.PrepTime += regex.Digit.FindString(matches[0]) + "H" + regex.Digit.FindString(matches[1]) + "M"
			}
		}

		if cook != "" {
			rs.CookTime = "PT" + regex.Digit.FindString(cook)
			if strings.Contains(cook, "min") {
				rs.CookTime += "M"
			} else {
				rs.CookTime += "H"
			}
		}
	}

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      regex.Digit.FindString(root.Find("li[title=Calories] span.top").Text()),
		Carbohydrates: regex.Digit.FindString(root.Find("li[title=Carbs] span.top").Text()),
		Fat:           regex.Digit.FindString(root.Find("li[title=Fat] span.top").Text()),
		Fiber:         regex.Digit.FindString(root.Find("li[title=Fibre] span.top").Text()),
		Protein:       regex.Digit.FindString(root.Find("li[title=Protein] span.top").Text()),
		SaturatedFat:  regex.Digit.FindString(root.Find("li[title=Saturates] span.top").Text()),
		Sodium:        regex.Digit.FindString(root.Find("li[title=Salt] span.top").Text()),
	}

	nodes := root.Find(".tags-list a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, sel.Text())
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	return rs, nil
}
