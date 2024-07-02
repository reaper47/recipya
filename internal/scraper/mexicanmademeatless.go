package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strings"
)

func scrapeMexicanMadeMeatless(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()
	rs.Description.Value = getNameContent(root, "description")
	rs.Image.Value = getPropertyContent(root, "og:image:secure_url")
	rs.DateModified = getPropertyContent(root, "og:updated_time")
	rs.Name = strings.TrimSpace(root.Find("h2.wprm-recipe-name ").Text())

	prep := strings.TrimSpace(regex.Digit.FindString(root.Find(".wprm-recipe-prep_time-minutes").Text()))
	if prep != "" {
		rs.PrepTime = "PT" + prep + "M"
	}

	cook := strings.TrimSpace(regex.Digit.FindString(root.Find(".wprm-recipe-cook_time-minutes").Text()))
	if cook != "" {
		rs.CookTime = "PT" + prep + "M"
	}

	xc := strings.Split(strings.TrimSpace(root.Find(".wprm-recipe-course").Text()), ",")
	if len(xc) > 0 {
		rs.Category.Value = strings.ToLower(strings.TrimSpace(xc[0]))
	}

	xc = strings.Split(strings.TrimSpace(root.Find(".wprm-recipe-cuisine").Text()), ",")
	if len(xc) > 0 {
		rs.Cuisine.Value = strings.ToLower(strings.TrimSpace(xc[0]))
	}

	rs.Yield.Value = findYield(root.Find(".wprm-recipe-servings-with-unit").Text())

	root.Find(".wprm-recipe-ingredient-group").Each(func(_ int, group *goquery.Selection) {
		group.Children().Each(func(_ int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "ul" {
				sel.Find("li").Each(func(_ int, sel *goquery.Selection) {
					s := strings.TrimSpace(sel.Text())
					if s != "" {
						rs.Ingredients.Values = append(rs.Ingredients.Values, s)
					}
				})
			} else {
				s := strings.TrimSpace(sel.Text())
				if s != "" {
					rs.Ingredients.Values = append(rs.Ingredients.Values, s)
				}
			}
		})
	})

	root.Find(".wprm-recipe-instruction-group").Each(func(_ int, group *goquery.Selection) {
		var name string
		group.Children().Each(func(_ int, sel *goquery.Selection) {
			if goquery.NodeName(sel) == "h4" {
				name = strings.TrimSpace(sel.Text())
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(name, &models.HowToItem{Name: name}))
			} else {
				sel.Find("li").Each(func(_ int, li *goquery.Selection) {
					s := strings.TrimSpace(li.Text())
					rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s, &models.HowToItem{Name: name}))

				})
			}
		})
	})

	notes := strings.TrimSpace(root.Find(".wprm-recipe-notes-container ").Text())
	if notes != "" {
		notes = strings.Join(strings.Fields(notes), " ")
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(notes, &models.HowToItem{Name: "notes"}))
	}

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      root.Find(".wprm-nutrition-label-text-nutrition-container-calories").Text(),
		Carbohydrates: root.Find(".wprm-nutrition-label-text-nutrition-container-carbohydrates").Text(),
		Fat:           root.Find(".wprm-nutrition-label-text-nutrition-container-fat").Text(),
		Fiber:         root.Find(".wprm-nutrition-label-text-nutrition-container-calories").Text(),
		Protein:       root.Find(".wprm-nutrition-label-text-nutrition-container-protein").Text(),
		SaturatedFat:  root.Find(".wprm-nutrition-label-text-nutrition-container-saturated_fat").Text(),
		Servings:      root.Find(".wprm-nutrition-label-text-nutrition-container-serving_size").Text(),
	}

	return rs, nil
}
