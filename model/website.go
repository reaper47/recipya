package model

// Websites holds a list of URLs.
type Websites struct {
	Objects []string `json:"websites"`
}

// Categories holds the URL of a website.
type Website struct {
	Url string `json:"url"`
}
