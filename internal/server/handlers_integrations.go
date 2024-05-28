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
		var (
			err         error
			integration = r.FormValue("integration")
			rawURL      = r.FormValue("url")
			username    = r.FormValue("username")
			password    = r.FormValue("password")
		)

		go func(id int64) {
			userIDAttr := slog.Int64("userID", id)

			s.Brokers.SendProgressStatus("Contacting server...", true, 0, -1, id)

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
				s.Brokers.HideNotification(id)
				s.Brokers.SendToast(models.NewErrorGeneralToast(msg), id)
				return
			case <-progress:
				for p := range progress {
					processed++
					s.Brokers.SendProgress("Fetching recipes...", processed, p.Total*2, id)
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
				s.Brokers.SendProgress("Adding to collection...", p.Value+p.Total, p.Total*2, id)
			}

			var (
				count   = len(recipeIDs)
				skipped = len(recipes) - count
			)

			slog.Info("Imported recipes", "integration", integration, userIDAttr, "count", count, "skipped", skipped)
			s.Brokers.HideNotification(id)
			s.Brokers.SendToast(models.NewInfoToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "", ""), id)
		}(getUserID(r))

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
			s.Brokers.SendToast(models.NewErrorGeneralToast("Connection failed. Please verify credentials."), getUserID(r))
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			s.Brokers.SendToast(models.NewInfoToast("Connection successful", msg+" "+"server connection verified.", ""), getUserID(r))
		}
	}
}
