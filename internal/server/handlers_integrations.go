package server

import (
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"log/slog"
	"net/http"
	"slices"
)

func (s *Server) integrationsImport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		_, found := s.Brokers[userID]
		if !found {
			w.Header().Set("HX-Trigger", models.NewWarningWSToast("Connection lost. Please reload page.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var (
			err         error
			integration = r.FormValue("integration")
			rawURL      = r.FormValue("url")
			username    = r.FormValue("username")
			password    = r.FormValue("password")
		)

		go func(id int64) {
			s.Brokers[id].SendProgressStatus("Contacting server...", true, 0, -1)

			var (
				recipes   models.Recipes
				processed int
				progress  = make(chan models.Progress)
				errs      = make(chan error, 1)
			)

			go func() {
				defer close(progress)

				switch integration {
				case "mealie":
					recipes, err = s.Integrations.MealieImport(rawURL, username, password, s.Files, progress)
				case "nextcloud":
					recipes, err = s.Integrations.NextcloudImport(rawURL, username, password, s.Files, progress)
				case "tandoor":
					recipes, err = s.Integrations.TandoorImport(rawURL, username, password, s.Files, progress)
				default:
					err = errors.New("no integration selected")
				}

				if err != nil {
					errs <- err
				} else if recipes == nil {
					errs <- errors.New("recipes from " + integration + " is nil")
				}
			}()

			select {
			case err = <-errs:
				msg := "Failed to import " + integration + " recipes."
				if integration == "" {
					msg = "No integration selected."
				}

				slog.Error(msg, userIDAttr, "processed", processed, "error", err)
				s.Brokers[id].HideNotification()
				s.Brokers[id].SendToast(models.NewErrorGeneralToast(msg))
				return
			case <-progress:
				for p := range progress {
					processed++
					s.Brokers[id].SendProgress("Fetching recipes...", processed, p.Total*2)
					close(progress)
				}
			}

			recipes = slices.DeleteFunc(recipes, func(r models.Recipe) bool {
				return len(r.Instructions) == 0 && len(r.Ingredients) == 0 && r.Description == ""
			})

			var (
				progress2 = make(chan models.Progress)
				recipeIDs []int64
			)

			go func() {
				defer close(progress2)

				recipeIDs, _, err = s.Repository.AddRecipes(recipes, id, progress2)
				if err != nil {
					slog.Error("Failed to add recipes", userIDAttr, "error", err)
				}
			}()

			for p := range progress2 {
				s.Brokers[id].SendProgress("Adding to collection...", p.Value+p.Total, p.Total*2)
			}

			var (
				count   = len(recipeIDs)
				skipped = len(recipes) - count
			)

			slog.Info("Imported recipes", "integration", integration, userIDAttr, "count", count, "skipped", skipped)
			s.Brokers[id].HideNotification()
			s.Brokers[id].SendToast(models.NewInfoToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "", ""))
		}(userID)

		w.WriteHeader(http.StatusAccepted)
	}
}

// TestConnection tests the connection of an integration. No error is returned on success.
func (s *Server) integrationTestConnectionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			api = r.URL.Query().Get("api")
			err = s.Integrations.TestConnection(api)
			msg string
		)

		switch api {
		case "azure-di":
			msg = "Azure AI Document Intelligence"
		case "sg":
			msg = "Twilio SendGrid"
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Connection failed. Please verify credentials.").Render())
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.Header().Set("HX-Trigger", models.NewInfoToast("Connection successful", msg+" "+"server connection verified.", "").Render())
		}
	}
}
