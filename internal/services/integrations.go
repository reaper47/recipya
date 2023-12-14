package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
)

// NewIntegrationsService creates a new Integrations that satisfies the IntegrationsService interface.
func NewIntegrationsService() *Integrations {
	return &Integrations{}
}

// Integrations is the entity that manages software integrations.
type Integrations struct{}

// NextcloudImport imports the recipes from a Nextcloud instance.
func (i *Integrations) NextcloudImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (*models.Recipes, error) {
	return integrations.NextcloudImport(baseURL, username, password, files.UploadImage, progress)
}

// ProcessImageOCR processes an image using an OCR service to extract the recipe.
func (i *Integrations) ProcessImageOCR(file io.Reader) (models.Recipe, error) {
	body := &bytes.Buffer{}
	_, err := io.Copy(body, file)
	if err != nil {
		return models.Recipe{}, err
	}

	url := app.Config.Integrations.AzureComputerVision.VisionEndpoint + "/computervision/imageanalysis:analyze?features=caption,read&model-version=latest&language=en&api-version=2023-02-01-preview"
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return models.Recipe{}, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", app.Config.Integrations.AzureComputerVision.ResourceKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.Recipe{}, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var av models.AzureVision
	err = json.NewDecoder(res.Body).Decode(&av)
	if err != nil {
		return models.Recipe{}, err
	}

	recipe := av.Recipe()
	if recipe.IsEmpty() {
		return recipe, errors.New("recipe is empty")
	}
	return recipe, nil
}
