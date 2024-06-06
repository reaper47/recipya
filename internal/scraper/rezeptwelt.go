package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeRezeptwelt(root *goquery.Document) (models.RecipeSchema, error) {
	rs.Name, _ = root.Find("meta[itemprop='name']").Attr("content")
	category, _ := root.Find("span[itemprop='recipeCategory']").Attr("content")
	rs.Description.Value, _ = root.Find("meta[itemprop='description']").Attr("content")
	rs.Image.Value, _ = root.Find("img[itemprop='image']").Attr("src")
	rs.DatePublished, _ = root.Find("meta[itemprop='datePublished']").Attr("content")
	rs.DateModified, _ = root.Find("meta[itemprop='dateModified']").Attr("content")

	nodes := root.Find("li[itemprop='recipeIngredient']")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.TrimSpace(s.Text())
	})

	prepTime, _ := root.Find("meta[itemprop='performTime']").Attr("content")
	cuisine, _ := root.Find("meta[itemprop='recipeCuisine']").Attr("content")
	keywords, _ := root.Find("meta[itemprop='keywords']").Attr("content")

	nodes = root.Find("meta[itemprop='tool']")
	tools := make([]string, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		t, _ := s.Attr("content")
		tools = append(tools, t)
	})

	nodes = root.Find("ol[itemprop='recipeInstructions'] li")
	instructions := make([]models.HowToStep, 0, nodes.Length())
	nodes.Each(func(_ int, s *goquery.Selection) {
		instructions = append(instructions, models.NewHowToStep(s.Text()))
	})

	yield := findYield(root.Find("span[itemprop='recipeYield']").Text())

	return models.RecipeSchema{
		Name:          name,
		Category:      &models.Category{Value: category},
		Cuisine:       &models.Cuisine{Value: cuisine},
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   &models.Description{Value: description},
		Image:         &models.Image{Value: image},
		Keywords:      &models.Keywords{Values: keywords},
		PrepTime:      prepTime,
		Yield:         &models.Yield{Value: yield},
		Ingredients:   &models.Ingredients{Values: ingredients},
		Instructions:  &models.Instructions{Values: instructions},
	}, nil
}
