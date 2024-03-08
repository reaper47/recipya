package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"net/http"
	"net/url"
	"strings"
)

const atContext = "https://schema.org"

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// IScraper is the scraper's interface.
type IScraper interface {
	Scrape(url string, files services.FilesService) (models.RecipeSchema, error)
}

// Scraper represents the IScraper's implementation.
type Scraper struct {
	Client HTTPClient
}

// NewScraper creates a new Scraper.
func NewScraper(client HTTPClient) *Scraper {
	return &Scraper{
		Client: client,
	}
}

// Scrape extracts the recipe from the given URL. An error will be
// returned when the URL cannot be parsed.
func (s *Scraper) Scrape(url string, files services.FilesService) (models.RecipeSchema, error) {
	host := getHost(url)
	if host == "bergamot" {
		return s.scrapeBergamot(url)
	} else if host == "foodbag" {
		return s.scrapeFoodbag(url)
	}

	doc, err := s.fetchDocument(url)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs, err := scrapeWebsite(doc, host)
	if err != nil {
		return rs, ErrNotImplemented
	}

	if rs.AtContext == "" {
		rs.AtContext = atContext
	}

	if rs.URL == "" {
		rs.URL = url
	}

	imageUUID, err := files.ScrapeAndStoreImage(rs.Image.Value)
	if err != nil {
		return rs, err
	}

	rs.Image.Value = imageUUID.String()
	return rs, nil
}

func (s *Scraper) fetchDocument(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	const mozilla = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0"

	switch getHost(url) {
	case "aberlehome", "thepalatablelife":
		req.Header.Set("User-Agent", mozilla)
	case "ah":
		req.Header.Set("Accept-Language", "nl")
		req.Header.Set("User-Agent", mozilla)
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("could not fetch (%d), try the bookmarklet", res.StatusCode)
	}

	return goquery.NewDocumentFromReader(res.Body)
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
		if s == "recipes" || s == "receitas" || s == "cooking" || s == "news" || s == "mobile" || s == "dashboard" || s == "fr" {
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
