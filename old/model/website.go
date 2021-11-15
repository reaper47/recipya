package model

import "github.com/reaper47/recipya/data"

// Websites holds a list of URLs.
type Websites struct {
	Objects [len(data.Websites)]string `json:"websites"`
}

// Categories holds the URL of a website.
type Website struct {
	Url string `json:"url"`
}
