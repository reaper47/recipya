package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

// Scrape extracts the recipe from the given URL. An error will be
// returned when the URL cannot be parsed.
func Scrape(rawurl string) (rs models.RecipeSchema, err error) {
	res, err := http.Get(rawurl)
	if err != nil {
		return rs, fmt.Errorf("could not fetch the url: %s", err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		return rs, fmt.Errorf("could not parse HTML: %s", err)
	}

	rs, err = scrapeWebsite(doc, getHost(rawurl))
	if rs.Url == "" {
		rs.Url = rawurl
	}
	return rs, nil
}

func getHost(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}

	parts := strings.Split(u.Hostname(), ".")
	switch len(parts) {
	case 4:
		if parts[1] == "m" {
			return parts[2]
		}
		return parts[1]
	case 3:
		if parts[0] == "recipes" || parts[0] == "receitas" || parts[0] == "cooking" {
			return parts[1]
		}

		if parts[1] == "wikibooks" || parts[1] == "tesco" {
			return parts[1]
		}

		if parts[0] != "www" {
			return parts[0]
		}

		return parts[1]
	default:
		return parts[len(parts)-2]

	}
}

func scrapeLdJSONs(root *html.Node) (rs models.RecipeSchema, err error) {
	n := getElement(root, "type", "application/ld+json")

	var xrs []models.RecipeSchema
	err = json.Unmarshal([]byte(n.FirstChild.Data), &xrs)
	if err != nil {
		return rs, err
	}

	for _, rs := range xrs {
		if rs.AtType.Value == "Recipe" {
			return rs, nil
		}
	}
	return models.RecipeSchema{}, nil
}

func scrapeLdJSON(root *html.Node) (rs models.RecipeSchema, err error) {
	n := getElement(root, "type", "application/ld+json")

	err = json.Unmarshal([]byte(n.FirstChild.Data), &rs)
	if err != nil {
		return rs, err
	}

	if rs.AtType.Value != "Recipe" {
		return rs, fmt.Errorf("@type must be Recipe but got %s", rs.AtType.Value)
	}
	return rs, nil
}

type graph struct {
	AtContext string                `json:"@context"`
	AtGraph   []models.RecipeSchema `json:"@graph"`
}

func scrapeGraph(root *html.Node) (rs models.RecipeSchema, err error) {
	n := getElement(root, "type", "application/ld+json")

	var g graph
	err = json.Unmarshal([]byte(n.FirstChild.Data), &g)
	if err != nil {
		return rs, err
	}

	for _, r := range g.AtGraph {
		if r.AtType.Value == "Recipe" {
			return r, nil
		}
	}
	return rs, fmt.Errorf("no recipe for the given url")
}

func findRecipeLdJSON(root *html.Node) (rs models.RecipeSchema, err error) {
	xn := traverseAll(root, func(node *html.Node) bool {
		return getAttr(node, "type") == "application/ld+json"
	})
	for _, n := range xn {
		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, "\n", "")

		r, err := scrapeLdJSON(n)
		if err != nil {
			continue
		}

		if r.AtType.Value == "Recipe" {
			return r, nil
		}
	}
	return rs, err
}
