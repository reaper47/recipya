package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapeSmittenKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	yield := int16(1)
	parsed, err := strconv.ParseInt(regex.Digit.FindString(root.Find("*[itemprop='recipeYield']").Text()), 10, 16)
	if err == nil {
		yield = int16(parsed)
	}

	prepTime, _ := root.Find("time[itemprop='totalTime']").Attr("datetime")

	nodes := root.Find(".jetpack-recipe-ingredients").Last().Find("ul")
	children := nodes.Children()
	ingredients := make([]string, 0, children.Length())
	children.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find(".jetpack-recipe-directions").Last()
	children = nodes.Children()
	instructions := make([]string, 0, children.Length())
	children.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		if s != "" {
			instructions = append(instructions, s)
		}
	})

	img, _ := root.Find(".post-thumbnail-container").First().Find("img").Attr("src")
	datePub, _ := root.Find(".entry-date.published").First().Attr("datetime")
	dateMod, _ := root.Find(".updated").First().Attr("datetime")

	title := root.Find("h3[itemprop='name']").Text()
	if title == "" {
		title = strings.TrimSpace(root.Find(".entry-title").First().Text())
	}

	if len(ingredients) == 0 && len(instructions) == 0 {
		root.Find(".smittenkitchen-print-hide").First().NextAll().Each(func(_ int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "p" {
				s := strings.TrimSpace(sel.Text())
				if strings.HasPrefix(strings.ToLower(s), "serve") {
					parsed, err = strconv.ParseInt(regex.Digit.FindString(s), 10, 16)
					if err == nil {
						yield = int16(parsed)
					}
					return
				}

				idx := regex.Unit.FindStringIndex(s)
				isAtStart := idx != nil && idx[0] < 8

				idx = regex.Digit.FindStringIndex(s)
				isStartWithNumber := idx != nil && idx[0] == 0

				if len(ingredients) == 0 && (isAtStart || isStartWithNumber) && sel.Has("br").Length() > 0 {
					xs := strings.Split(sel.Text(), "\n")
					ingredients = make([]string, 0, len(xs))
					for _, x := range xs {
						ingredients = append(ingredients, strings.TrimSpace(x))
					}
				} else if len(ingredients) > 0 {
					instructions = append(instructions, strings.TrimSpace(sel.Text()))
				}
			}
		})
	}

	return models.RecipeSchema{
		DateModified:    dateMod,
		DatePublished:   datePub,
		Description:     models.Description{Value: root.Find(".smittenkitchen-print-hide p").First().Text()},
		Image:           models.Image{Value: img},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            title,
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        prepTime,
		Yield:           models.Yield{Value: yield},
	}, nil
}
