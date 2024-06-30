package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeGrandfrais(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value = getPropertyContent(root, "og:description")
	rs.Name = root.Find("h1[itemprop='name']").Text()
	rs.Yield.Value = findYield(root.Find("p[itemprop='recipeYield']").Text())

	prep := root.Find(".pre-requie-item p").First().Text()
	split := strings.Split(prep, " ")
	isMin := strings.Contains(prep, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			prep = "PT" + split[i] + "M"
		}
	}
	rs.PrepTime = prep

	cook := root.Find(".pre-requie-item p").Last().Text()
	split = strings.Split(cook, " ")
	isMin = strings.Contains(cook, "min")
	for i, s := range split {
		_, err := strconv.ParseInt(s, 10, 64)
		if err == nil && isMin {
			cook = "PT" + split[i] + "M"
		}
	}
	rs.CookTime = cook

	nodes := root.Find("div[itemprop='ingredients']")
	rs.Ingredients.Values = strings.Split(nodes.Text(), "\n")
	for i, ingredient := range rs.Ingredients.Values {
		_, after, ok := strings.Cut(ingredient, "- ")
		if ok {
			rs.Ingredients.Values[i] = after
		}
	}

	getInstructions(&rs, root.Find("div[itemprop='recipeInstructions'] li"))

	image, _ := root.Find("img[itemprop='image']").Attr("src")
	if !strings.HasPrefix(image, "https://") {
		image = "https://www.grandfrais.com" + image
	}
	rs.Image.Value = image

	return rs, nil
}
