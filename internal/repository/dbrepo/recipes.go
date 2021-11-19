package db

import "github.com/reaper47/recipya/internal/models"

// GetAllRecipes gets all of the recipes in the database.
func (m *postgresDBRepo) GetAllRecipes() ([]models.Recipe, error) {
	recipes := []models.Recipe{}
	return recipes, nil
}
