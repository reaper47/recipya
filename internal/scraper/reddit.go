package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
	"strings"
)

func scrapeReddit(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	node := root.Find(".commentarea").First()

	datePub, _ := node.Find("time").First().Attr("datetime")
	rs.DatePublished, _, _ = strings.Cut(datePub, "T")
	rs.Name = root.Find("a[data-event-action='title']").Text()

	form := node.Find(".sitetable.nestedlisting form").First()
	nodes := form.Find("ul").First().Find("li")
	rs.Ingredients.Values = make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Ingredients.Values = append(rs.Ingredients.Values, sel.Text())
	})

	node = form.Find("p:contains('Instructions')")
	if node.Nodes != nil {
		for c := node.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				s := strings.TrimSpace(c.Data)
				dotIndex := strings.IndexByte(s, '.')
				if dotIndex < 4 {
					_, s, _ = strings.Cut(s, ".")
				}
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(s))
			}
		}
	}

	return rs, nil
}
