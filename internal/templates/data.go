package templates

// Data holds data to pass on to the templates.
type Data struct {
	IsAuthenticated bool // IsAuthenticated says whether the user is authenticated.

	Title string // Title is the text inserted <title> tag's text.

	Content      string // Content is text to insert into the template.
	ContentTitle string // ContentTitle is the header of the Content.

	Scraper ScraperData
}

// RegisterData is the data to pass on to the user registration template.
type RegisterData struct {
	Email           string
	PasswordConfirm string
}

// ScraperData holds template data related to the recipe scraper.
type ScraperData struct {
	UnsupportedWebsite string
}
