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
	srv, ts, c := createWSServer()
	defer func() {
		_ = c.Close()
	}()

	originalRepo := srv.Repository
	originalIntegrations := srv.Integrations

	uriImport := ts.URL + "/integrations/import/nextcloud"

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
			assertHeader(t, rr, "HX-Trigger", `{"showToast":"{\"background\":\"alert-error\",\"message\":\"Invalid username, password or URL.\",\"title\":\"\"}"}`)
		})
	}

	t.Run("error when importing", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			NextcloudImportFunc: func(baseURL, username, password string, files services.FilesService) (*models.Recipes, error) {
				return nil, errors.New("import error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"{\"background\":\"alert-error\",\"message\":\"Failed to import Nextcloud recipes.\",\"title\":\"\"}"}`
		assertWebsocket(t, c, 3, want)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered:      make(map[int64]models.Recipes),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		srv.Integrations = &mockIntegrations{
			NextcloudImportFunc: func(baseURL, username, password string, files services.FilesService) (*models.Recipes, error) {
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

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"{\"background\":\"alert-info\",\"message\":\"Imported 2 recipes. Skipped 0.\",\"title\":\"\"}"}`
		assertWebsocket(t, c, 5, want)
		if len(repo.RecipesRegistered[1]) != 2 {
			t.Fatal("expected 2 recipes in the repo")
		}
	})
}
