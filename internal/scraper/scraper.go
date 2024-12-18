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
	"net"
	"net/http"
	"strings"
	"time"
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
func (s *Scraper) Scrape(rawURL string, files services.FilesService) (models.RecipeSchema, error) {
	var host = s.HTTP.GetHost(rawURL)
	rs, isSpecial, err := s.scrapeSpecial(host, rawURL, files)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	if isSpecial {
		img, _ := files.ScrapeAndStoreImage(rs.Image.Value)
		rs.Image.Value = img.String()
		return rs, nil
	}

	if host == "reddit" {
		rawURL = strings.Replace(rawURL, "www", "old", 1)
	}

	doc, err := s.fetchDocument(rawURL)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	rs, err = scrapeWebsite(doc, host)
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
		rs.URL = rawURL
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

func (s *Scraper) scrapeSpecial(host, rawURL string, files services.FilesService) (models.RecipeSchema, bool, error) {
	var (
		rs        models.RecipeSchema
		err       error
		isSpecial = true
	)

	switch host {
	case "bergamot":
		rs, err = s.scrapeBergamot(rawURL)
	case "foodbag":
		rs, err = s.scrapeFoodbag(rawURL)
	case "gousto":
		rs, err = s.scrapeGousto(rawURL)
	case "madewithlau":
		rs, err = s.scrapeMadeWithLau(rawURL)
	case "monsieur-cuisine":
		doc, err := s.fetchDocument(rawURL)
		if err != nil {
			return models.RecipeSchema{}, true, err
		}
		rs, err = s.scrapeMonsieurCuisine(doc, rawURL, files)
	default:
		isSpecial = false
	}

	return rs, isSpecial, err
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

	// Some websites require cookies
	host := s.HTTP.GetHost(url)
	switch host {
	case "bettybossi":
		_ = res.Body.Close()
		res, err = s.HTTP.Client.Do(req)
		if err != nil {
			return nil, err
		}
	case "chatelaine", "damndelicious", "dinnerthendessert":
		original := s.HTTP.Client
		c, ok := s.HTTP.Client.(*http.Client)
		if ok {
			c.Transport = &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 10 * time.Second,
			}

			res, err = s.HTTP.Client.Do(req)
			if err != nil {
				return nil, err
			}

			s.HTTP.Client = original
		}
	}

	defer res.Body.Close()

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
