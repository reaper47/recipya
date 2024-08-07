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
	rs.DatePublished, _, _ = strings.Cut(node.Find("time").First().AttrOr("datetime", ""), "T")
	rs.Name = root.Find("a[data-event-action='title']").Text()

	form := node.Find(".sitetable.nestedlisting form").First()
	getIngredients(&rs, form.Find("ul").First().Find("li"))

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
