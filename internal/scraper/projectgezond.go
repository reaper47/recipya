package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeProjectgezond(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, found := strings.Cut(name, " | ")
	if found {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	rs.Description.Value, _ = root.Find("meta[property='og:description']").Attr("content")
	rs.Category.Value, _ = root.Find("meta[property='article:section']").Attr("content")
	rs.Image.Value, _ = root.Find(".wp-post-image").First().Attr("src")

	rs.DatePublished, _ = root.Find("meta[property='og:published_time']").Attr("content")
	rs.DateModified, _ = root.Find("meta[property='article:modified_time']").Attr("content")

	nodes := root.Find("h2").First().NextUntil("h2")
	ingredientNodes := nodes.Find("ul li")
	rs.Ingredients.Values = make([]string, 0, ingredientNodes.Length())
	ingredientNodes.Each(func(_ int, sel *goquery.Selection) {
		s := sel.Text()
		rs.Ingredients.Values = append(rs.Ingredients.Values, s)
	})

	nodes = nodes.Next().NextUntil("h2")
	instructionNodes := nodes.Find("ul li")
	rs.Instructions.Values = make([]models.HowToItem, 0, instructionNodes.Length())
	instructionNodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	var cal string
	node := root.Find("strong:contains('Kcal')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		cal = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var protein string
	node = root.Find("strong:contains('Eiwit')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		protein = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var carbs string
	node = root.Find("strong:contains('Koolhydraten')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		carbs = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var fat string
	node = root.Find("strong:contains('Vet')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		fat = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	var fiber string
	node = root.Find("strong:contains('Vezels')")
	if len(node.Nodes) > 0 && node.Nodes[0].NextSibling != nil {
		fiber = strings.TrimSpace(node.Nodes[0].NextSibling.Data)
	}

	rs.NutritionSchema = &models.NutritionSchema{
		Calories:      cal + " kcal",
		Carbohydrates: carbs,
		Fat:           fat,
		Fiber:         fiber,
		Protein:       protein,
	}

	return rs, nil
}
