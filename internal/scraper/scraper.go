package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const atContext = "https://schema.org"

// Scrape extracts the recipe from the given URL. An error will be
// returned when the URL cannot be parsed.
func Scrape(url string) (models.RecipeSchema, error) {
	var rs models.RecipeSchema

	res, err := http.Get(url)
	if err != nil {
		return rs, fmt.Errorf("could not fetch the url: %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return rs, fmt.Errorf("could not parse HTML: %s", err)
	}

	rs, err = scrapeWebsite(doc, getHost(url))
	if err != nil {
		return rs, errors.New("url could not be fetched")
	}

	if rs.AtContext == "" {
		rs.AtContext = atContext
	}

	if rs.URL == "" {
		rs.URL = url
	}

	if rs.Image.Value != "" {
		res, err := http.Get(rs.Image.Value)
		if err != nil {
			return rs, nil
		}
		defer res.Body.Close()

		img, _, err := image.Decode(res.Body)
		if err != nil {
			return rs, nil
		}

		imageUUID := uuid.New().String()
		out, err := os.Create(filepath.Join(app.ImagesDir, imageUUID+".jpg"))
		if err != nil {
			return rs, nil
		}
		defer out.Close()

		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		if width > 800 || height > 800 {
			img = imaging.Resize(img, width/2, height/2, imaging.NearestNeighbor)
		}

		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 33}); err != nil {
			return rs, nil
		}

		rs.Image.Value = imageUUID
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
		if parts[0] == "recipes" || parts[0] == "receitas" || parts[0] == "cooking" {
			return parts[1]
		}

		if parts[1] == "wikibooks" || parts[1] == "tesco" || parts[1] == "expressen" {
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

func parseLdJSON(root *goquery.Document) (models.RecipeSchema, error) {
	for _, node := range root.Find("script[type='application/ld+json']").Nodes {
		var rs models.RecipeSchema
		if err := json.Unmarshal([]byte(node.FirstChild.Data), &rs); err != nil {
			var xrs []models.RecipeSchema
			if err := json.Unmarshal([]byte(node.FirstChild.Data), &xrs); err != nil {
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
	return models.RecipeSchema{}, nil
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
		if err := json.Unmarshal([]byte(node.FirstChild.Data), &g); err != nil {
			continue
		}

		for _, r := range g.AtGraph {
			if r.AtType.Value == "Recipe" {
				r.AtContext = atContext
				return r, nil
			}
		}
	}
	return models.RecipeSchema{}, fmt.Errorf("no recipe for the given url")
}
