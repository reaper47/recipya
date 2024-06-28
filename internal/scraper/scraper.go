package scraper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"io"
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
		HTTP: services.NewHTTP(client),
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

	// Some websites require cookies
	host := s.HTTP.GetHost(url)
	switch host {
	case "bettybossi", "chatelaine":
		res, err = s.HTTP.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("could not fetch (%d), try the bookmarklet", res.StatusCode)
	}

	var r io.Reader = res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		zr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		defer zr.Close()

		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, zr)
		if err != nil {
			return nil, err
		}

		r = buf
	}

	return goquery.NewDocumentFromReader(r)
}
