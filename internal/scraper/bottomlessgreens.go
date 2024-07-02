package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"slices"
	"strings"
)

func scrapeBottomLessGreens(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Name = getItempropContent(root, "headline")
	rs.Description.Value = getItempropContent(root, "description")
	rs.ThumbnailURL.Value = getItempropContent(root, "thumbnailUrl")
	rs.Image.Value = getItempropContent(root, "image")
	rs.DatePublished = getItempropContent(root, "datePublished")
	rs.DateModified = getItempropContent(root, "dateModified")

	nodes := root.Find(".blog-meta-item--tags").Last().Find("a")
	xk := make([]string, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		xk = append(xk, strings.TrimSpace(sel.Text()))
	})
	rs.Keywords.Values = strings.Join(xk, ",")

	u := getItempropContent(root, "url")
	_, after, ok := strings.Cut(u, "https://bottomlessgreens.com/")
	if ok {
		before, _, _ := strings.Cut(after, "/")
		rs.Category.Value = strings.ReplaceAll(before, "-", " ")
	}

	getIngredients(&rs, root.Find("ul").First().Find("li"))

	root.Find("h3").Each(func(_ int, sel *goquery.Selection) {
		name := strings.TrimSpace(sel.Text())

		for c := sel.Nodes[0].NextSibling; c != nil; c = c.NextSibling {
			if c.Data == "h3" {
				break
			}

			text := strings.TrimSpace(goquery.NewDocumentFromNode(c).Text())
			if text == "" {
				continue
			}

			if !slices.ContainsFunc(rs.Instructions.Values, func(item models.HowToItem) bool {
				return item.Text == text
			}) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(text, &models.HowToItem{
					Name: name,
				}))
			}
		}
	})

	tools := root.Find("h4:contains('Kitchen “stuff”')").First().Next().Text()
	parts := strings.Split(tools, ",")
	rs.Tools.Values = make([]models.HowToItem, 0, len(parts))
	for _, part := range parts {
		s := strings.TrimSuffix(strings.TrimSpace(part), ".")
		rs.Tools.Values = append(rs.Tools.Values, models.NewHowToTool(s))
	}

	return rs, nil
}
