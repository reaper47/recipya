package templates

import "github.com/reaper47/recipya/internal/models"

// IndexData holds data to pass on to the index template.
type RecipesData struct {
	Recipes []models.Recipe
}

// RecipeData holds data to pass to the recipe templates.
type RecipeData struct {
	Recipe models.Recipe
}
