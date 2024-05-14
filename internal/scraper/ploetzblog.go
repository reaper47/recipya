package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

func scrapePloetzblog(root *goquery.Document) (models.RecipeSchema, error) {
	name, _ := root.Find("meta[property='og:title']").Attr("content")
	image, _ := root.Find("img.we2p-pb-recipe__thumbnail-image").Attr("src")

	var description strings.Builder
	root.Find("div.we2p-pb-recipe__description p").Each(func(_ int, sel *goquery.Selection) {
		description.WriteString(strings.TrimSpace(sel.Text()))
		description.WriteString("\n\n")
	})

	nodes := root.Find(".we2p-pb-recipe__ingredients table").Last().Find("tr")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		parts := strings.Fields(sel.Text())
		if len(parts) > 3 {
			parts = parts[:3]
		}
		s := strings.Join(parts, " ")
		ingredients = append(ingredients, strings.TrimSpace(s))
	})

	yield, _ := root.Find("#recipePieceCount").Attr("value")

	total := strings.ToLower(root.Find("span:contains('Gesamtzubereitungszeit')").Parent().Children().Last().Text())
	var prep string
	before, after, ok := strings.Cut(total, "stunden")
	if ok {
		prep = "PT" + regex.Digit.FindString(before) + "H" + regex.Digit.FindString(after) + "M"
	}

	var instructions []string
	nodes = root.Find("h4")
	nodes.Each(func(_ int, sel *goquery.Selection) {
		n := sel.Next().Children().Last()
		if n == nil || len(n.Nodes) == 0 || n.Nodes[0].Data != "div" {
			return
		}

		var sb strings.Builder
		sb.WriteString(strings.TrimSpace(sel.Text()) + "\n")
		n.Find("p").Each(func(i int, subSel *goquery.Selection) {
			sb.WriteString(strconv.Itoa(i+1) + ". " + strings.Join(strings.Fields(subSel.Text()), " ") + "\n")
		})
		instructions = append(instructions, sb.String()+"\n")
	})

	return models.RecipeSchema{
		Description:  models.Description{Value: strings.TrimSpace(description.String())},
		Image:        models.Image{Value: image},
		Ingredients:  models.Ingredients{Values: ingredients},
		Instructions: models.Instructions{Values: instructions},
		Name:         name,
		PrepTime:     prep,
		Yield:        models.Yield{Value: findYield(yield)},
	}, nil
}
