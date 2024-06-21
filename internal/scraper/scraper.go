package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"net/http"
	"strings"
)

const atContext = "https://schema.org"

// IScraper is the scraper's interface.
type IScraper interface {
	Scrape(url string, files services.FilesService) (models.RecipeSchema, error)
}

// Scraper represents the IScraper's implementation.
type Scraper struct {
	HTTP services.HTTP
}

// NewScraper creates a new Scraper.
func NewScraper(client services.HTTPClient) *Scraper {
	return &Scraper{
		HTTP: *services.NewHTTP(client),
	}
}

// Scrape extracts the recipe from the given URL. An error will be
// returned when the URL cannot be parsed.
func (s *Scraper) Scrape(url string, files services.FilesService) (models.RecipeSchema, error) {
	host := s.HTTP.GetHost(url)
	if host == "bergamot" {
		return s.scrapeBergamot(url)
	} else if host == "foodbag" {
		return s.scrapeFoodbag(url)
	} else if host == "monsieur-cuisine" {
		doc, err := s.fetchDocument(url)
		if err != nil {
			return models.RecipeSchema{}, err
		}
		return s.scrapeMonsieurCuisine(doc, url, files)
	} else if host == "quitoque" {
		return s.scrapeQuitoque(url)
	} else if host == "reddit" {
		url = strings.Replace(url, "www", "old", 1)
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

	if rs.Category == nil || rs.Category.Value == "" {
		rs.Category = models.NewCategory("")
	}

	if rs.AtType == nil {
		rs.AtType = &models.SchemaType{Value: "Recipe"}
	} else if rs.AtType.Value == "" {
		rs.AtType.Value = "Recipe"
	}

	if rs.URL == "" {
		rs.URL = url
	}

	var imageUUID uuid.UUID
	if rs.Image != nil {
		imageUUID, err = files.ScrapeAndStoreImage(rs.Image.Value)
		if err != nil {
			return rs, err
		}
		rs.Image.Value = imageUUID.String()
	}

	return rs, nil
}

func (s *Scraper) fetchDocument(url string) (*goquery.Document, error) {
	req, err := s.HTTP.PrepareRequestForURL(url)
	if err != nil {
		return nil, err
	}

	res, err := s.HTTP.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Some websites require page reloads
	host := s.HTTP.GetHost(url)
	if host == "bettybossi" {
		res, err = s.HTTP.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("could not fetch (%d), try the bookmarklet", res.StatusCode)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func parseLdJSON(root *goquery.Document) (models.RecipeSchema, error) {
	for _, node := range root.Find("script[type='application/ld+json']").Nodes {
		if node.FirstChild == nil {
			continue
		}

		var rs = models.NewRecipeSchema()
		err := json.Unmarshal([]byte(strings.ReplaceAll(node.FirstChild.Data, "\n", "")), &rs)
		if err != nil || rs.Equal(models.NewRecipeSchema()) {
			var xrs []models.RecipeSchema
			err = json.Unmarshal([]byte(node.FirstChild.Data), &xrs)
			if err != nil {
				continue
			}

			for _, rs = range xrs {
				if rs.AtType != nil && rs.AtType.Value == "Recipe" {
					return rs, nil
				}
			}
			continue
		}

		if rs.AtType.Value != "Recipe" {
			continue
		}
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
	return models.NewRecipeSchema(), ErrNotImplemented
}
