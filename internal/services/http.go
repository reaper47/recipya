package services

import (
	"net/http"
	"net/url"
	"strings"
)

func GetGetRequestForUrl(url string) (*http.Request, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	const mozilla = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.0"

	host := GetHost(url)
	switch host {
	case "aberlehome", "bettybossi", "marmiton", "natashaskitchen", "puurgezond", "reddit", "thepalatablelife":
		req.Header.Set("User-Agent", mozilla)
	case "ah":
		req.Header.Set("Accept-Language", "nl")
		req.Header.Set("User-Agent", mozilla)
	}

	return req, err
}

func GetHost(rawURL string) string {
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
