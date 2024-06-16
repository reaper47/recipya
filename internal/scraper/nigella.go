package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"regexp"
	"strconv"
	"strings"
)

func scrapeNigella(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	name, _ := root.Find("meta[property='og:title']").Attr("content")
	before, _, ok := strings.Cut(name, " — ")
	if ok {
		name = strings.TrimSpace(before)
	}
	rs.Name = name

	yieldNode := root.Find("p[class='serves']")
	yieldText := strings.Split(yieldNode.Nodes[0].FirstChild.Data, ":")[1]
	numberPattern, _ := regexp.Compile("[0-9]+")
	yieldNumMatch := numberPattern.FindString(yieldText)
	yield, err := strconv.Atoi(yieldNumMatch)
	if err != nil {
		return rs, err
	}
	rs.Yield.Value = int16(yield)

	untrimmedDescription, _ := root.Find("meta[property='og:description']").Attr("content")
	description, _ := strings.CutSuffix(untrimmedDescription, "\n\nFor US cup measures, use the toggle at the top of the ingredients list.")
	rs.Description.Value = description
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")

	nodes := root.Find("li[itemprop=recipeIngredient]")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(sel.Text()))
	})

	nodes = root.Find("div[itemprop=recipeInstructions]").First().Find("ol li")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(sel.Text())))
	})

	return rs, nil
}
