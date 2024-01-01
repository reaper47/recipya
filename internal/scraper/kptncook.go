package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strconv"
	"strings"
)

func scrapeKptncook(root *goquery.Document) (models.RecipeSchema, error) {
	image, _ := root.Find("image[itemprop='image']").Attr("src")

	var prep string
	clock := root.Find("img[alt='clockIcon']")
	if clock != nil && len(clock.Nodes) > 0 && clock.Nodes[0].NextSibling != nil {
		data := clock.Nodes[0].NextSibling.Data
		isMin := strings.Contains(data, "min")
		split := strings.Split(data, " ")
		for i, s := range split {
			_, err := strconv.ParseInt(s, 10, 64)
			if err == nil && isMin {
				prep = "PT" + split[i] + "M"
				break
			}
		}
	}

	// It is possible to extract the category, but I am not sure how.

	var ingredients []string
	node := root.Find("img[src='/assets/images/img_wave_grey4.png']").Parent().Parent()
	for {
		node = node.Next()
		if node.Nodes == nil {
			break
		}
		s := strings.Join(strings.Fields(node.Text()), " ")
		ingredients = append(ingredients, s)
	}

	nodes := root.Find(".kptn-step-title")
	instructions := make([]string, 0, nodes.Length())
	nodes.Each(func(i int, sel *goquery.Selection) {
		s := strings.TrimSpace(sel.Text())
		_, err := strconv.ParseInt(s[:1], 10, 64)
		if err == nil {
			_, after, _ := strings.Cut(s, ".")
			after = strings.Join(strings.Fields(after), " ")
			s = strings.TrimSpace(after)
		}
		instructions = append(instructions, s)
	})

	return models.RecipeSchema{
		AtContext:       atContext,
		AtType:          models.SchemaType{Value: "Recipe"},
		Category:        models.Category{},
		CookingMethod:   models.CookingMethod{},
		Cuisine:         models.Cuisine{},
		DateCreated:     "",
		DateModified:    "",
		DatePublished:   "",
		Description:     models.Description{},
		Keywords:        models.Keywords{},
		Image:           models.Image{Value: image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            root.Find("title").Text(),
		NutritionSchema: models.NutritionSchema{},
		PrepTime:        prep,
		Tools:           models.Tools{},
		Yield:           models.Yield{Value: findYield(root.Find(".kptn-person-count").Text())},
	}, nil
}
