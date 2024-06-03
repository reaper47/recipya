package server_test

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	"net/http"
	"strings"
	"testing"
)

func TestHandlers_Integrations_Import(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	originalRepo := srv.Repository
	originalIntegrations := srv.Integrations

	uriImport := ts.URL + "/integrations/import"

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodPost, uriImport)
	})

	t.Run("missing integration", func(t *testing.T) {
		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"No integration selected.","title":"General Error"}}`
		assertWebsocket(t, c, 3, want)
	})

	t.Run("error when importing", func(t *testing.T) {
		srv.Integrations = &mockIntegrations{
			importFunc: func(baseURL, username, password string, files services.FilesService) (models.Recipes, error) {
				return nil, errors.New("import error")
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("integration=nextcloud&username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Failed to import nextcloud recipes.","title":"General Error"}}`
		assertWebsocket(t, c, 3, want)
	})

	t.Run("valid request", func(t *testing.T) {
		repo := &mockRepository{
			RecipesRegistered:      make(map[int64]models.Recipes),
			UserSettingsRegistered: map[int64]*models.UserSettings{1: {}},
		}
		srv.Repository = repo
		srv.Integrations = &mockIntegrations{
			importFunc: func(baseURL, username, password string, files services.FilesService) (models.Recipes, error) {
				return models.Recipes{
					{ID: 1, Name: "One", Ingredients: []string{"one"}},
					{ID: 2, Name: "Two", Ingredients: []string{"two"}},
				}, nil
			},
		}
		defer func() {
			srv.Integrations = originalIntegrations
			srv.Repository = originalRepo
		}()

		rr := sendHxRequestAsLoggedIn(srv, http.MethodPost, uriImport, formHeader, strings.NewReader("integration=nextcloud&username=admin&password=admin&url=http://localhost:8080"))

		assertStatus(t, rr.Code, http.StatusAccepted)
		want := `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-info","message":"","title":"Imported 2 recipes. Skipped 0."}}`
		assertWebsocket(t, c, 5, want)
		if len(repo.RecipesRegistered[1]) != 2 {
			t.Fatal("expected 2 recipes in the repo")
		}
	})
}

func TestHandlers_Integrations_TestConnection(t *testing.T) {
	srv, ts, c := createWSServer()
	defer c.CloseNow()

	original := srv.Integrations
	uri := ts.URL + "/integrations/test-connection"

	type testcase struct {
		name string
		api  string
	}

	t.Run("must be logged in", func(t *testing.T) {
		assertMustBeLoggedIn(t, srv, http.MethodGet, uri)
	})

	testcases := []testcase{
		{name: "empty", api: ""},
		{name: "random", api: "random"},
	}
	for _, tc := range testcases {
		t.Run("invalid api "+tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?api="+tc.api)

			assertStatus(t, rr.Code, http.StatusBadRequest)
		})
	}

	testcases2 := []testcase{
		{name: "azure document intelligence", api: "azure-di"},
		{name: "sendgrid", api: "sg"},
	}
	for _, tc := range testcases2 {
		t.Run("invalid credentials "+tc.name, func(t *testing.T) {
			srv.Integrations = &mockIntegrations{
				testConnectionFunc: func(api string) error {
					return errors.New("connection failed")
				},
			}
			defer func() {
				srv.Integrations = original
			}()
			rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?api="+tc.api)

			assertStatus(t, rr.Code, http.StatusUnauthorized)
			assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-error","message":"Connection failed. Please verify credentials.","title":"General Error"}}`)
		})
	}

	testcases3 := []struct {
		name  string
		api   string
		toast string
	}{
		{name: "azure document intelligence", api: "azure-di", toast: "Azure AI Document Intelligence server connection verified."},
		{name: "sendgrid", api: "sg", toast: "Twilio SendGrid server connection verified."},
	}
	for _, tc := range testcases3 {
		t.Run("valid credentials "+tc.name, func(t *testing.T) {
			rr := sendHxRequestAsLoggedInNoBody(srv, http.MethodGet, uri+"?api="+tc.api)

			assertStatus(t, rr.Code, http.StatusOK)
			assertWebsocket(t, c, 1, `{"type":"toast","fileName":"","data":"","toast":{"action":"","background":"alert-info","message":"`+tc.toast+`","title":"Connection successful"}}`)
		})
	}
}
