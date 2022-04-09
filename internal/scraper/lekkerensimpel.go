package scraper

import (
	"strconv"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeLekkerenSimpel(root *html.Node) (models.RecipeSchema, error) {
	chCategory := make(chan string)
	chYield := make(chan int16)
	chPrep := make(chan string)
	chCook := make(chan string)
	go func() {
		var (
			category string
			yield    int16
			prep     string
			cook     string
		)
		defer func() {
			_ = recover()
			chCategory <- category
			chYield <- yield
			chPrep <- prep
			chCook <- cook
		}()

		numLi := 1
		ul := getElement(root, "class", "recipe__meta").FirstChild.NextSibling
		for c := ul.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}
			node := getElement(c, "class", "recipe__meta-title")

			switch numLi {
			case 1:
				category = strings.TrimSpace(node.FirstChild.Data)
				category = strings.ReplaceAll(category, "\n", "")
			case 2:
				s := strings.TrimSpace(node.FirstChild.Data)
				s = strings.ReplaceAll(s, "\n", "")
				parts := strings.Split(s, " ")
				for _, part := range parts {
					i, err := strconv.Atoi(part)
					if err == nil {
						yield = int16(i)
					}
				}
			case 3:
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.TextNode {
						continue
					}

					s := strings.ToLower(strings.TrimSpace(c.Data))
					s = strings.ReplaceAll(s, "\n", "")

					parts := strings.Split(s, " ")
					var minutes int
					for _, part := range parts {
						i, err := strconv.Atoi(part)
						if err == nil {
							minutes = i
						}
					}

					hours, minutes := minutes/60, minutes%60
					time := "PT"
					if hours != 0 {
						time += strconv.Itoa(hours) + "H"
					}
					time += strconv.Itoa(minutes) + "M"

					if strings.Contains(s, "bereidingstijd") {
						prep = time
					} else if strings.Contains(s, "oventijd") {
						cook = time
					}
				}
			}
			numLi++
		}
	}()

	chDescription := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chDescription <- s
		}()

		var vals []string
		node := getElement(root, "class", "entry__content")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			if c.Data == "p" {
				if len(vals) > 0 {
					vals = append(vals, "\n\n")
				}

				for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
					switch c2.Type {
					case html.TextNode:
						vals = append(vals, c2.Data)
					case html.ElementNode:
						for c3 := c2.FirstChild; c3 != nil; c3 = c3.NextSibling {
							switch c3.Type {
							case html.TextNode:
								vals = append(vals, c3.Data)
							case html.ElementNode:
								vals = append(vals, c3.FirstChild.Data)
							}
						}
					}
				}
			}

			if c.Data == "h2" {
				break
			}
		}
		s = strings.Join(vals, "")
		s = strings.TrimSpace(s)
	}()

	chDatePublished := make(chan string)
	chDateModified := make(chan string)
	go func() {
		var pub string
		var mod string
		defer func() {
			_ = recover()
			chDatePublished <- pub
			chDateModified <- mod
		}()

		pub = getAttr(getElement(root, "property", "article:published_time"), "content")
		mod = getAttr(getElement(root, "property", "article:modified_time"), "content")
	}()

	chIngredients := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- vals
		}()

		xn := traverseAll(getElement(root, "class", "entry__content"), func(node *html.Node) bool {
			return strings.HasPrefix(strings.ToLower(node.Data), "benodigdheden")
		})
		if len(xn) == 0 {
			xn = traverseAll(getElement(root, "class", "recipe__sidebar-title h3"), func(node *html.Node) bool {
				return strings.HasPrefix(strings.ToLower(node.Data), "benodigdheden")
			})
			if len(xn) == 0 {
				return
			}
		}
		node := xn[0]

		for {
			for c := node; c != nil; c = c.NextSibling {
				if c.Type != html.ElementNode {
					node = c.Parent
					continue
				}

				if c.Data == "ul" {
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						if c.Type != html.ElementNode {
							continue
						}
						vals = append(vals, c.FirstChild.Data)
					}
					return
				}

				if c.Data == "div" {
					for c := c.FirstChild; c != nil; c = c.NextSibling {
						if c.Type != html.ElementNode {
							continue
						}

						if c.Data == "ul" {
							for c := c.FirstChild; c != nil; c = c.NextSibling {
								if c.Type != html.ElementNode {
									continue
								}

								value := strings.TrimSpace(<-getElementData(c, "class", "value"))
								measure := strings.TrimSpace(<-getElementData(c, "class", "measure"))
								name := strings.TrimSpace(<-getElementData(c, "class", "ingredient"))

								var s string
								if value != "" {
									s += value + " "
								}

								if measure != "" {
									s += measure + " "
								}

								if name != "" {
									s += name
								}
								vals = append(vals, s)
							}
							break
						}
					}
					return
				}
			}
			node = node.Parent
		}
	}()

	chInstructions := make(chan []string)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- vals
		}()

		xn := traverseAll(getElement(root, "class", "entry__content"), func(node *html.Node) bool {
			s := strings.ToLower(node.Data)
			return strings.HasPrefix(s, "bereidingswijze") || strings.HasPrefix(s, "aan de slag")
		})

		if len(xn) == 0 {
			return
		}

		for c := xn[0].Parent.NextSibling; c != nil; c = c.NextSibling {
			if c.Type != html.ElementNode {
				continue
			}

			switch c.Data {
			case "p":
				var xs []string
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					switch c.Type {
					case html.ElementNode:
						xs = append(xs, strings.TrimSpace(c.FirstChild.Data))
					case html.TextNode:
						xs = append(xs, strings.TrimSpace(c.Data))
					}
				}
				vals = append(vals, strings.Join(xs, " "))
			case "ol":
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.ElementNode {
						continue
					}
					vals = append(vals, strings.TrimSpace(c.FirstChild.Data))
				}
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:     "https://schema.org",
		AtType:        models.SchemaType{Value: "Recipe"},
		Name:          <-getElementData(root, "class", "hero__title"),
		Category:      models.Category{Value: <-chCategory},
		Yield:         models.Yield{Value: <-chYield},
		PrepTime:      <-chPrep,
		CookTime:      <-chCook,
		Description:   models.Description{Value: <-chDescription},
		Ingredients:   models.Ingredients{Values: <-chIngredients},
		Instructions:  models.Instructions{Values: <-chInstructions},
		DatePublished: <-chDatePublished,
		DateModified:  <-chDateModified,
	}, nil
}
