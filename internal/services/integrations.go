package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/integrations"
	"github.com/reaper47/recipya/internal/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
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

// TandoorImport imports recipes from a Tandoor instance.
func (i Integrations) TandoorImport(baseURL, username, password string, files FilesService, progress chan models.Progress) (models.Recipes, error) {
	if !isCredentialsValid(baseURL, username, password) {
		return nil, errors.New("invalid username, password or URL")
	}
	return integrations.TandoorImport(baseURL, username, password, i.client, files.UploadImage, progress)
}

func isCredentialsValid(baseURL, username, password string) bool {
	if username == "" || password == "" || baseURL == "" {
		return false
	}

	_, err := url.Parse(baseURL)
	return err == nil
}

// ProcessImageOCR processes an image using an OCR service to extract the recipe.
func (i Integrations) ProcessImageOCR(files []io.Reader) (models.Recipes, error) {
	locations := make([]string, 0, len(files))
	for _, file := range files {
		req, err := app.Config.Integrations.AzureDI.PrepareRequest(file)
		if err != nil {
			slog.Error("Failed to prepare request", "error", err)
			continue
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Failed to execute request", "req", req, "err", err)
			continue
		}

		if res.StatusCode == http.StatusAccepted {
			locations = append(locations, res.Header.Get("Operation-Location"))
		} else {
			slog.Warn("Failed to execute request", "req", req, "status", res.StatusCode)
		}

		_ = res.Body.Close()
	}

	var (
		recipes  = make(models.Recipes, 0, len(locations))
		idx      = 0
		maxTries = 0
	)

	for idx < len(locations) {
		locAttr := slog.String("loc", locations[idx])

		if maxTries > 10 {
			slog.Warn("Too many tries", locAttr, "maxTries", maxTries)
			idx++
			maxTries = 0
			continue
		}

		req, err := http.NewRequest(http.MethodGet, locations[idx], nil)
		if err != nil {
			slog.Warn("Failed to prepare request", "location", locations[idx], "err", err)
			idx++
			continue
		}
		req.Header.Set("Ocp-Apim-Subscription-Key", app.Config.Integrations.AzureDI.Key)

		res, err := i.client.Do(req)
		if err != nil {
			slog.Warn("Failed to execute request", locAttr, "err", err)
			idx++
			maxTries = 0
			continue
		}

		if res.StatusCode != http.StatusOK {
			var adi models.AzureDIError
			_ = json.NewDecoder(res.Body).Decode(&adi)
			slog.Warn("Undesirable status code", locAttr, "status", res.StatusCode, "data", adi)
			_ = res.Body.Close()
			idx++
			maxTries = 0
			continue
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Warn("Failed to read response body", locAttr, "err", err)
		}
		slog.Info("Processed Azure AI DI response", locAttr, "body", string(body))

		var adi models.AzureDILayout
		err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&adi)
		if err != nil {
			slog.Warn("Failed to decode response body", locAttr, "err", err)
			_ = res.Body.Close()
			idx++
			maxTries = 0
			continue
		}

		if adi.Status != "succeeded" {
			time.Sleep(1 * time.Second)
			maxTries++
			continue
		}

		recipe := adi.Recipe()
		if recipe.IsEmpty() {
			slog.Warn("OCR recipe is empty", locAttr, "azure response", adi)
		} else {
			recipes = append(recipes, recipe)
		}

		_ = res.Body.Close()
		idx++
		maxTries = 0
	}

	return recipes, nil
}

// TestConnection tests the connection of an integration. No error is returned on success.
func (i Integrations) TestConnection(api string) error {
	var (
		apiAttr       = slog.String("api", api)
		errConnFailed = errors.New("connection failed")
	)

	switch api {
	case "azure-di":
		req, err := app.Config.Integrations.AzureDI.PrepareRequest(nil)
		if err != nil {
			slog.Error("Failed to prepare request", apiAttr, "error", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Failed to send request", "req", req, apiAttr, "error", err)
			return err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusNotFound {
			return errConnFailed
		}
		return nil
	case "sg":
		client := sendgrid.NewSendClient(app.Config.Email.SendGridAPIKey)
		res, err := client.Send(mail.NewSingleEmail(mail.NewEmail("", ""), "", mail.NewEmail("", ""), "", ""))
		if err != nil {
			slog.Error("Failed to send request", apiAttr, "error", err)
			return err
		}

		if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusNotFound {
			return errConnFailed
		}
		return nil
	default:
		return errors.New("invalid api")
	}
}
