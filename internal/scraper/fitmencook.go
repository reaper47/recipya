package scraper

import (
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeFitMenCook(root *html.Node) (rs models.RecipeSchema, err error) {
	content := getElement(root, "id", "fit-wrapper")

	chKeywords := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chKeywords <- vals
		}()

		node := getElement(content, "class", "fit-post-categories")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				v := c.FirstChild.FirstChild.Data
				vals = append(vals, strings.ToLower(v))
			}
		}
	}()

	chPrepTime := make(chan string)
	chCookTime := make(chan string)
	go func() {
		var prep, cook string
		defer func() {
			_ = recover()
			chPrepTime <- prep
			chCookTime <- cook
		}()

		str := getElement(content, "class", "prep-time macros").FirstChild.Data
		if str != "" {
			parts := strings.Split(str, " ")
			if strings.ToLower(parts[1]) == "minutes" {
				prep = "PT" + parts[0] + "M"
			}

		}

		str = getElement(content, "class", "cook-time macros").FirstChild.Data
		if str != "" {
			parts := strings.Split(str, " ")
			if strings.ToLower(parts[1]) == "minutes" {
				cook = "PT" + parts[0] + "M"
			}

		}
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var v models.Ingredients
		defer func() {
			_ = recover()
			chIngredients <- v
		}()

		node := getElement(content, "class", "recipe-ingredients gap-bottom-small")
		var ul *html.Node
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data != "ul" {
				continue
			}
			ul = c
			break
		}

		for l := ul.FirstChild; l != nil; l = l.NextSibling {
			if l.Type != html.ElementNode {
				continue
			}

			for l2 := l.FirstChild; l2 != nil; l2 = l2.NextSibling {
				if l2.Type != html.ElementNode {
					continue
				}

				switch l2.Data {
				case "strong":
					v.Values = append(v.Values, "\n", l2.FirstChild.Data)
				case "ul":
					for ing := l2.FirstChild; ing != nil; ing = ing.NextSibling {
						if ing.Type != html.ElementNode {
							continue
						}

						var s string
						for el := ing.FirstChild; el != nil; el = el.NextSibling {
							switch el.Type {
							case html.TextNode:
								s += el.Data
							case html.ElementNode:
								s += el.FirstChild.Data
							}
						}
						v.Values = append(v.Values, s)
					}
				}
			}
		}
		v.Values = v.Values[1:]
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var v models.Instructions
		defer func() {
			_ = recover()
			chInstructions <- v
		}()

		node := getElement(content, "class", "recipe-steps")
		var ol *html.Node
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data != "ol" {
				continue
			}
			ol = c
			break
		}

		for li := ol.FirstChild; li != nil; li = li.NextSibling {
			if li.Type != html.ElementNode {
				continue
			}
			s := li.FirstChild.Data
			s = strings.ReplaceAll(s, "\u00a0", "")
			v.Values = append(v.Values, s)
		}
	}()

	rs, err = scrapeLdJSON(root)
	rs.Keywords = models.Keywords{Values: strings.TrimSpace(strings.Join(<-chKeywords, ","))}
	rs.PrepTime = <-chPrepTime
	rs.CookTime = <-chCookTime
	rs.Ingredients = <-chIngredients
	rs.Instructions = <-chInstructions
	return rs, err
}
