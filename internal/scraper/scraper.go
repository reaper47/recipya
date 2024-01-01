package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"net/http"
	"net/url"
	"strings"
)

const atContext = "https://schema.org"

// Scrape extracts the recipe from the given URL. An error will be
// returned when the URL cannot be parsed.
func Scrape(url string, files services.FilesService) (models.RecipeSchema, error) {
	var rs models.RecipeSchema

	res, err := http.Get(url)
	if err != nil {
		return rs, fmt.Errorf("could not fetch the url: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res == nil {
		return rs, errors.New("scrape response is nil")
	}

	body := res.Body
	if body == nil {
		return rs, errors.New("scrape response body is nil")
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return rs, fmt.Errorf("could not parse HTML: %w", err)
	}

	rs, err = scrapeWebsite(doc, getHost(url))
	if err != nil {
		return rs, ErrNotImplemented
	}

	if rs.AtContext == "" {
		rs.AtContext = atContext
	}

	if rs.URL == "" {
		rs.URL = url
	}

	if rs.Image.Value != "" {
		resImage, err := http.Get(rs.Image.Value)
		if err != nil {
			return rs, err
		}
		defer func() {
			_ = resImage.Body.Close()
		}()

		if resImage == nil {
			return rs, errors.New("image response is nil")
		}

		body = resImage.Body
		if body == nil {
			return rs, errors.New("image response body is nil")
		}

		imageUUID, err := files.UploadImage(body)
		if err != nil {
			return rs, err
		}
		rs.Image.Value = imageUUID.String()
	}
	return rs, nil
}

func getHost(rawURL string) string {
	u, err := url.Parse(rawURL)
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
		s := parts[0]
		if s == "recipes" || s == "receitas" || s == "cooking" || s == "news" || s == "mobile" {
			return parts[1]
		}

		if parts[1] == "wikibooks" || parts[1] == "tesco" || parts[1] == "expressen" {
			return parts[1]
		}

		if s != "www" {
			return s
		}
		return parts[1]
	default:
		if len(parts) > 1 {
			return parts[len(parts)-2]
		}
		return ""
	}
}

func parseLdJSON(root *goquery.Document) (models.RecipeSchema, error) {
	for _, node := range root.Find("script[type='application/ld+json']").Nodes {
		if node.FirstChild == nil {
			continue
		}

		var rs models.RecipeSchema
		err := json.Unmarshal([]byte(strings.ReplaceAll(node.FirstChild.Data, "\n", "")), &rs)
		if err != nil {
			var xrs []models.RecipeSchema
			err = json.Unmarshal([]byte(node.FirstChild.Data), &xrs)
			if err != nil {
				continue
			}

			for _, rs := range xrs {
				if rs.AtType.Value == "Recipe" {
					rs.AtContext = atContext
					return rs, nil
				}
			}
			continue
		}

		if rs.AtType.Value != "Recipe" {
			continue
		}

		rs.AtContext = atContext
		return rs, nil
	}
	return models.RecipeSchema{}, ErrNotImplemented
}

type graph struct {
	AtContext string                `json:"@context"`
	AtGraph   []models.RecipeSchema `json:"@graph"`
}

func parseGraph(root *goquery.Document) (models.RecipeSchema, error) {
	for _, node := range root.Find("script[type='application/ld+json']").Nodes {
		if node.FirstChild == nil {
			continue
		}

		var g graph
		err := json.Unmarshal([]byte(node.FirstChild.Data), &g)
		if err != nil {
			continue
		}

		for _, r := range g.AtGraph {
			if r.AtType.Value == "Recipe" {
				r.AtContext = atContext
				return r, nil
			}
		}
	}
	return models.RecipeSchema{}, ErrNotImplemented
}
