package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeArchanasKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	description := root.Find("span[itemprop='description']").Text()
	description = strings.TrimPrefix(description, "\n")
	description = strings.ReplaceAll(description, "\u00a0", " ")
	description = strings.TrimSpace(description)

	image, _ := root.Find("img[itemprop='image']").Attr("src")
	image = "https://www.archanaskitchen.com" + image

	var keywords string
	root.Find("li[itemprop='keywords'] a").Each(func(_ int, s *goquery.Selection) {
		keywords += strings.TrimSpace(s.Text()) + ","
	})
	keywords = strings.TrimSuffix(keywords, ",")

	nodes := root.Find("li[itemprop='ingredients']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\n", "")
		v = strings.ReplaceAll(v, "\t", "")
		ingredients[i] = strings.TrimSpace(strings.ReplaceAll(v, " , ", ", "))
	})

	nodes = root.Find("li[itemprop='recipeInstructions'] p")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		v := strings.ReplaceAll(s.Text(), "\u00a0", " ")
		v = strings.TrimSpace(strings.ReplaceAll(v, " .", "."))
		instructions[i] = v
	})

	prepTime, _ := root.Find("span[itemprop='prepTime']").Attr("content")
	cookTime, _ := root.Find("span[itemprop='cookTime']").Attr("content")
	datePublished, _ := root.Find("span[itemprop='datePublished']").Attr("content")
	dateModified, _ := root.Find("span[itemprop='dateModified']").Attr("content")
	yield := findYield(root.Find("span[itemprop='recipeYield'] p").Text())

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          root.Find("h1[itemprop='name']").Text(),
		Description:   models.Description{Value: description},
		Image:         models.Image{Value: image},
		Category:      models.Category{Value: root.Find(".recipeCategory a").Text()},
		Cuisine:       models.Cuisine{Value: root.Find("span[itemprop='recipeCuisine'] a").Text()},
		PrepTime:      prepTime,
		CookTime:      cookTime,
		DatePublished: datePublished,
		DateModified:  dateModified,
		Keywords:      models.Keywords{Values: keywords},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Yield:         models.Yield{Value: yield},
	}, nil
}
