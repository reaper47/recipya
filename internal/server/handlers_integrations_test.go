package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_Integrations_Nextcloud(t *testing.T) {
	srv := newServerTest()
	originalRepo := srv.Repository
	originalIntegrations := srv.Integrations

	uriImport := "/integrations/import/nextcloud"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uriImport)
	})

	invalidData := []struct {
		name string
		in   string
	}{
		{name: "no base URL", in: "username=admin&password=admin"},
		{name: "no username", in: "password=admin&url=http://localhost:8080"},
		{name: "no password", in: "username=admin&url=http://localhost:8080"},
	}
	for _, tc := range invalidData {
		t.Run(tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader(tc.in))

			assertStatus(t, rr.Code, http.StatusBadRequest)
			assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Invalid username, password or URL.\",\"backgroundColor\":\"bg-red-500\"}"}`)
		})
	}

	t.Run("error when importing", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			NextcloudImportFunc: func(client *http.Client, baseURL, username, password string, files services.FilesService) (*models.Recipes, error) {
				return nil, errors.New("import error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusInternalServerError)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Failed to import Nextcloud recipes.\",\"backgroundColor\":\"bg-red-500\"}"}`)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered: make(map[int64]models.Recipes),
		}
		srv.Repository = repo
		srv.Integrations = &mockIntegrations{
			NextcloudImportFunc: func(client *http.Client, baseURL, username, password string, files services.FilesService) (*models.Recipes, error) {
				return &models.Recipes{
					{ID: 1, Name: "One"},
					{ID: 2, Name: "Two"},
				}, nil
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusCreated)
		assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"message\":\"Imported 2 recipes. Skipped 0.\",\"backgroundColor\":\"bg-blue-500\"}"}`)
		if len(repo.RecipesRegistered[1]) != 2 {
			t.Fatal("expected 2 recipes in the repo")
		}
	})
}
