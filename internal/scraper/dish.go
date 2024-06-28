package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

func scrapeDish(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	description := getNameContent(root, "description")
	before, after, ok := strings.Cut(description, "PREP:")
	if ok {
		rs.Description.Value = before
		// PREP: 8 minutes&nbsp;&nbsp;COOK: 35 minutes&nbsp;&nbsp;SCALE: easy
		parts := strings.Split(after, "&nbsp;")
		atoi, err := strconv.Atoi(regex.Digit.FindString(parts[0]))
		if err == nil {
			if !strings.Contains(parts[0], "h") {
				rs.PrepTime = "PT" + strconv.Itoa(atoi) + "M"
			}
		}

		if len(parts) > 1 {
			for _, part := range parts[1:] {
				part = strings.ToLower(part)
				if strings.HasPrefix(part, "cook") {
					atoi, err = strconv.Atoi(regex.Digit.FindString(part))
					if err == nil {
						if !strings.Contains(parts[0], "h") {
							rs.CookTime = "PT" + strconv.Itoa(atoi) + "M"
						}
					}
				}
			}
		}
	} else {
		rs.Description.Value = description
	}

	rs.DatePublished = getPropertyContent(root, "article:published_time")
	rs.Name = strings.TrimSpace(root.Find("h1[itemprop='name']").Text())
	img, _ := root.Find("img[itemprop='image']").First().Attr("src")
	rs.Image.Value = "https://dish.co.nz/" + img
	rs.Yield.Value = findYield(root.Find("h2.serve").First().Text())

	root.Find("div[itemprop='ingredients']").First().Find("p").Each(func(_ int, node *goquery.Selection) {
		for c := node.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
			var v string
			if c.Type == html.TextNode {
				v = c.Data
			} else if c.Data == "strong" {
				v = c.FirstChild.Data
			}

			v = strings.TrimSpace(v)
			if v != "" {
				rs.Ingredients.Values = append(rs.Ingredients.Values, v)
			}
		}
	})

	nodes := root.Find("div[itemprop='recipeInstructions']").First().Find("p")
	rs.Instructions.Values = make([]models.HowToItem, 0, nodes.Length())
	nodes.Each(func(_ int, sel *goquery.Selection) {
		rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(sel.Text()))
	})

	return rs, nil
}
