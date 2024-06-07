package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/utils/regex"
	"golang.org/x/net/html"
	"slices"
	"strings"
)

func scrapePuurgezond(root *goquery.Document) (models.RecipeSchema, error) {
	rs := models.NewRecipeSchema()

	rs.Description.Value, _ = root.Find("meta[name='description']").Attr("content")
	rs.Keywords.Values, _ = root.Find("meta[name='keywords']").Attr("content")
	rs.Image.Value, _ = root.Find("meta[property='og:image']").Attr("content")
	rs.Name = root.Find("title").Text()

	var (
		isIngredients  bool
		isInstructions bool
	)

	root.Find("div[itemprop='articleBody']").Children().Each(func(_ int, sel *goquery.Selection) {
		s := strings.ToLower(sel.Text())
		if strings.Contains(s, "person") {
			for c := sel.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
				if c.Type != html.TextNode {
					continue
				}

				if strings.Contains(c.Data, "person") {
					rs.Yield.Value = findYield(c.Data)
				} else if strings.Contains(c.Data, "bereidingstijd") {
					match := slices.DeleteFunc(regex.Time.FindStringSubmatch(c.Data), func(s string) bool { return s == "" })
					switch len(match) {
					case 2:
						if strings.Contains(match[1], "h") || strings.Contains(match[1], "time") {
							rs.PrepTime = "PT" + regex.Digit.FindString(match[1]) + "H"
						} else if strings.Contains(match[1], "min") {
							rs.PrepTime = "PT" + regex.Digit.FindString(match[1]) + "M"
						}
					case 3:
						rs.PrepTime = "PT" + regex.Digit.FindString(match[1]) + "H" + regex.Digit.FindString(match[2]) + "M"
					}
				} else if strings.Contains(c.Data, "oventijd") {
					match := slices.DeleteFunc(regex.Time.FindStringSubmatch(c.Data), func(s string) bool { return s == "" })
					switch len(match) {
					case 2:
						if strings.Contains(match[1], "h") || strings.Contains(match[1], "time") {
							rs.CookTime = "PT" + regex.Digit.FindString(match[1]) + "H"
						} else if strings.Contains(match[1], "min") {
							rs.CookTime = "PT" + regex.Digit.FindString(match[1]) + "M"
						}
					case 3:
						rs.CookTime = "PT" + regex.Digit.FindString(match[1]) + "H" + regex.Digit.FindString(match[2]) + "M"
					}
				}
			}
		} else if strings.HasPrefix(s, "ingrediënten") {
			isIngredients = true
		} else if strings.HasPrefix(s, "bereiden") {
			isInstructions = true
		} else if isIngredients {
			for c := sel.Nodes[0].FirstChild; c != nil; c = c.NextSibling {
				if regex.Digit.FindStringIndex(c.Data) == nil && len(c.Data) > 50 {
					break
				} else if c.Type == html.TextNode {
					rs.Ingredients.Values = append(rs.Ingredients.Values, strings.TrimSpace(c.Data))
				}
			}
			isIngredients = false
		} else if isInstructions {
			sel.Find("li").Each(func(_ int, li *goquery.Selection) {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(strings.TrimSpace(li.Text())))
			})
			isInstructions = false
		}
	})

	return rs, nil
}
