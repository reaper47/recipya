package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
	"strings"
)

func scrapeReddit(root *goquery.Document) (models.RecipeSchema, error) {
	node := root.Find(".commentarea").First()

	datePub, _ := node.Find("time").First().Attr("datetime")
	datePub, _, _ = strings.Cut(datePub, "T")

	form := node.Find(".sitetable.nestedlisting form").First()
	nodes := form.Find("ul").First().Find("li")
	ingredients := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		ingredients = append(ingredients, sel.Text())
	})

	var instructions []string
	node = form.Find("p:contains('Instructions')")
	if node.Nodes != nil {
		for c := node.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				s := strings.TrimSpace(c.Data)
				dotIndex := strings.IndexByte(s, '.')
				if dotIndex < 4 {
					_, s, _ = strings.Cut(s, ".")
				}
				instructions = append(instructions, strings.TrimSpace(s))
			}
		}
	}

	return models.RecipeSchema{
		AtContext:     atContext,
		AtType:        models.SchemaType{Value: "Recipe"},
		DatePublished: datePub,
		Ingredients:   models.Ingredients{Values: ingredients},
		Instructions:  models.Instructions{Values: instructions},
		Name:          root.Find("a[data-event-action='title']").Text(),
		Yield:         models.Yield{},
	}, nil
}
