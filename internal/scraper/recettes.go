package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeRecettesDuQuebec(root *goquery.Document) (models.RecipeSchema, error) {
	rs.DatePublished, _ = root.Find("meta[name='cXenseParse:recs:publishtime']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Name, _ = root.Find("meta[name='cXenseParse:title']").Attr("content")
	rs.Image.Value, _ = root.Find("picture img").Attr("srcset")

	category := root.Find(".categories h6").Siblings().First().Text()
	category = strings.TrimSpace(category)

	keywords := root.Find(".tags h6").Siblings().First().Text()
	keywords = strings.TrimSpace(keywords)

	var (
		prepTime string
		cookTime string
		yield    int16
	)
	root.Find("span.cat").Each(func(_ int, sel *goquery.Selection) {
		switch sel.Text() {
		case "Préparation":
			prepTime = "PT"
			parts := strings.Split(sel.Siblings().First().Text(), "&")
			if len(parts) == 2 {
				prepTime += strings.Split(parts[0], " ")[0] + "H"
				prepTime += strings.Split(parts[1], " ")[1] + "M"
			} else {
				xs := strings.Split(parts[0], " ")
				prepTime += xs[0]
				if strings.Contains(xs[1], "min") {
					prepTime += "M"
				} else {
					prepTime += "H"
				}
			}
		case "Cuisson":
			cookTime = "PT"
			parts := strings.Split(sel.Siblings().First().Text(), "&")
			if len(parts) == 2 {
				cookTime += strings.Split(parts[0], " ")[0] + "H"
				cookTime += strings.Split(parts[1], " ")[1] + "M"
			} else {
				xs := strings.Split(parts[0], " ")
				cookTime += xs[0]
				if strings.Contains(xs[1], "min") {
					cookTime += "M"
				} else {
					cookTime += "H"
				}
			}
		case "Portion(s)":
			yieldStr := sel.Siblings().First().Text()
			for _, s := range strings.Split(yieldStr, " ") {
				yield64, err := strconv.ParseInt(s, 10, 16)
				if err == nil {
					yield = int16(yield64)
					break
				}
			}
		}
	})

	nodes := root.Find(".ingredients ul label")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		s = strings.ReplaceAll(s, "\n", "")
		s = strings.Join(strings.Fields(s), " ")
		ingredients = append(ingredients, s)
	})

	nodes = root.Find(".method p")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(sel.Text()))
	})

	var recipeImage string
	split := strings.Split(image, "?")
	if len(split) > 0 {
		recipeImage = split[0]
	}

	return models.RecipeSchema{
		Category:      &models.Category{Value: category},
		CookTime:      cookTime,
		DatePublished: datePublished,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: recipeImage},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
		Keywords:      &models.Keywords{Values: keywords},
		Name:          name,
		PrepTime:      prepTime,
		Yield:         &models.Yield{Value: yield},
	}, nil
}
