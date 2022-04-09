package scraper

import (
	"github.com/reaper47/recipya/internal/models"
	"golang.org/x/net/html"
)

func scrapeNourishedByNutrition(root *html.Node) (models.RecipeSchema, error) {
	rs, err := scrapeGraph(root)
	rs.CookTime = ""
	return rs, err
}
