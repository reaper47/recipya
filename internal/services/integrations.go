package services

import (
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
)

// NewIntegrationsService creates a new Integrations that satisfies the IntegrationsService interface.
func NewIntegrationsService() *Integrations {
	return &Integrations{}
}

// Integrations is the entity that manages software integrations.
type Integrations struct{}

func (i *Integrations) NextcloudImport(client *http.Client, baseURL, username, password string, files FilesService) (*models.Recipes, error) {
	return integrations.NextcloudImport(client, baseURL, username, password, files.UploadImage)
}
