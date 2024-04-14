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
	"net/url"
)

// NewIntegrationsService creates a new Integrations that satisfies the IntegrationsService interface.
func NewIntegrationsService(client *http.Client) Integrations {
	return Integrations{
		client: client,
	}
}

// Integrations is the entity that manages software integrations.
type Integrations struct {
	client *http.Client
}

// MealieImport imports the recipes from a Mealie instance.
func (i Integrations) MealieImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error) {
	if !isCredentialsValid(baseURL, username, password) {
		return nil, errors.New("invalid username, password or URL")
	}
	return integrations.MealieImport(baseURL, username, password, i.client, files.UploadImage, progress)
}

// NextcloudImport imports the recipes from a Nextcloud instance.
func (i Integrations) NextcloudImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error) {
	if !isCredentialsValid(baseURL, username, password) {
		return nil, errors.New("invalid username, password or URL")
	}
	return integrations.NextcloudImport(baseURL, username, password, files.UploadImage, progress)
}

func isCredentialsValid(baseURL, username, password string) bool {
	if username == "" || password == "" || baseURL == "" {
		return false
	}

	_, err := url.Parse(baseURL)
	return err == nil
}

// ProcessImageOCR processes an image using an OCR service to extract the recipe.
func (i Integrations) ProcessImageOCR(file io.Reader) (models.Recipe, error) {
	body := &bytes.Buffer{}
	_, err := io.Copy(body, file)
	if err != nil {
		return models.Recipe{}, err
	}

	url := app.Config.Integrations.AzureComputerVision.VisionEndpoint + "/computervision/imageanalysis:analyze?features=caption,read&model-version=latest&language=en&api-version=2024-02-01"
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
	defer res.Body.Close()

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
