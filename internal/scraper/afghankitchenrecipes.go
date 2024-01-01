package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeAfghanKitchen(root *goquery.Document) (models.RecipeSchema, error) {
	content := root.Find("#content")
	info := content.Find(".recipe-info")
	about := content.Find("div[itemprop='about']")
	yield := findYield(info.Find(".servings .value").Text())

	var prep string
	time := info.Find(".prep-time .value").Text()
	if strings.Contains(time, "m") {
		prep = "PT" + strings.TrimSuffix(time, "m") + "M"
	} else if strings.Contains(time, "h") {
		time = strings.TrimSuffix(time, " h")
		parts := strings.Split(time, ":")
		if len(parts) >= 2 {
			prep = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}

	var cook string
	time = info.Find(".cook-time .value").Text()
	if strings.Contains(time, "m") {
		cook = "PT" + strings.TrimSuffix(time, "m") + "M"
	} else if strings.Contains(time, "h") {
		time = strings.TrimSuffix(time, " h")
		parts := strings.Split(time, ":")
		if len(parts) > 1 {
			cook = "PT" + parts[0] + "H" + parts[1] + "M"
		}
	}

	var description string
	if len(about.Nodes) > 0 && about.Nodes[0].FirstChild != nil && about.Nodes[0].FirstChild.NextSibling != nil && about.Nodes[0].FirstChild.NextSibling.FirstChild != nil {
		description = about.Nodes[0].FirstChild.NextSibling.FirstChild.Data
		description = strings.ReplaceAll(description, "\n", "")
		description = strings.ReplaceAll(description, "\u00a0", " ")
	}

	nodes := about.Find("li.ingredient")
	ingredients := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		ingredients[i] = strings.ReplaceAll(s.Text(), "  ", " ")
	})

	nodes = about.Find("p.instructions")
	instructions := make([]string, nodes.Length())
	nodes.Each(func(i int, s *goquery.Selection) {
		instructions[i] = strings.ReplaceAll(strings.TrimSpace(s.Text()), "  ", " ")
	})

	datePub, _ := content.Find("meta[itemprop='datePublished']").Attr("content")
	image, _ := content.Find("meta[itemprop='image']").Attr("content")

	return models.RecipeSchema{
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          content.Find("h2[itemprop='name']").Text(),
		DatePublished: datePub,
		Image:         models.Image{Value: image},
		Yield:         models.Yield{Value: yield},
		PrepTime:      prep,
		CookTime:      cook,
		Description:   models.Description{Value: strings.TrimSpace(description)},
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Category:      models.Category{Value: about.Find(".type a").Text()},
	}, nil
}
