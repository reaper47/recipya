package models

// NextcloudRecipes holds a Nextcloud recipe's metadata obtained from the Nextcloud Cookbook's /recipes endpoint.
type NextcloudRecipes struct {
	Category            string `json:"category"`
	DateCreated         string `json:"dateCreated"`
	DateModified        string `json:"dateModified"`
	ID                  int64  `json:"recipe_id"`
	ImagePlaceholderURL string `json:"imagePlaceholderUrl"`
	ImageURL            string `json:"imageUrl"`
	Keywords            string `json:"keywords"`
}
