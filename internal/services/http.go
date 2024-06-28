package services

import (
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient is an interface for making HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTP creates a new HTTP service.
func NewHTTP(client HTTPClient) HTTP {
	if client == nil {
		client = http.DefaultClient
	}

	return HTTP{
		Client: client,
	}
}

// HTTP is the entity that manages the HTTP client.
type HTTP struct {
	Client HTTPClient
}

// Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
func (h HTTP) Do(req *http.Request) (*http.Response, error) {
	return h.Client.Do(req)
}

// PrepareRequestForURL Prepares an HTTP GET request for a given URL.
// It will apply additional HTTP headers if the host requires it.
func (h HTTP) PrepareRequestForURL(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	const mozilla = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0"

	host := h.GetHost(url)
	switch host {
	case "aberlehome", "bettybossi", "downshiftology", "findingtimeforcooking", "marmiton", "natashaskitchen",
		"puurgezond", "reddit", "sarahsveganguide", "thekitchn", "thepalatablelife", "wellplated":
		req.Header.Set("User-Agent", mozilla)
	case "ah":
		req.Header.Set("Accept-Language", "nl")
		req.Header.Set("User-Agent", mozilla)
	case "chatelaine":
		req.Header.Set("User-Agent", mozilla)
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("DNT", "1")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept", "*/*")
		req.Header.Add("Accept-Charset", "utf-8")
		req.Header.Set("Connection", "keep-alive")
	}

	return req, err
}

// GetHost gets the host from the raw URL.
func (h HTTP) GetHost(rawURL string) string {
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
		if s == "recipes" || s == "receitas" || s == "cooking" || s == "news" || s == "mobile" ||
			s == "dashboard" || s == "fr" || s == "blog" || s == "old" {
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
