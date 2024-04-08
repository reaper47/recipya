package server

import (
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"log/slog"
	"net/http"
)

func (s *Server) integrationsImportNextcloud() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		_, found := s.Brokers[userID]
		if !found {
			w.Header().Set("HX-Trigger", models.NewWarningWSToast("Connection lost. Please reload page.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			msg := "Failed to get user settings."
			slog.Error(msg, userIDAttr, "error", err)
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast(models.NewErrorDBToast(msg))
			return
		}

		go func(id int64, us models.UserSettings) {
			s.Brokers[id].SendProgressStatus("Contacting server...", true, 0, -1)

			var (
				recipes   *models.Recipes
				processed int
				progress  = make(chan models.Progress)
				errs      = make(chan error, 1)
			)

			go func() {
				defer close(progress)
				recipes, err = s.Integrations.NextcloudImport(r.FormValue("url"), r.FormValue("username"), r.FormValue("password"), s.Files, progress)
				if err != nil {
					errs <- err
				}

				if recipes == nil {
					errs <- errors.New("recipes from Nextcloud is nil")
				}
			}()

			select {
			case err = <-errs:
				msg := "Failed to import Nextcloud recipes."
				slog.Error(msg, userIDAttr, "processed", processed, "error", err)
				s.Brokers[id].HideNotification()
				s.Brokers[id].SendToast(models.NewErrorGeneralToast(msg))
				return
			case <-progress:
				for p := range progress {
					processed++
					s.Brokers[id].SendProgress("Fetching recipes...", processed, p.Total*2)
				}
			}

			count := 0
			skipped := 0
			numRecipes := len(*recipes)
			recipeIDs := make([]int64, 0, numRecipes)

			for i, recipe := range *recipes {
				s.Brokers[id].SendProgress("Adding to collection...", i+numRecipes, numRecipes*2)
				c := recipe.Copy()
				recipeID, err := s.Repository.AddRecipe(&c, id, us)
				if err != nil {
					slog.Warn("Skipped recipe", userIDAttr, "recipe", c, "error", err)
					skipped++
					continue
				}
				recipeIDs = append(recipeIDs, recipeID)
				count++
			}

			slog.Info("Imported Nextcloud recipes", userIDAttr, "count", count, "skipped", skipped)
			s.Repository.CalculateNutrition(userID, recipeIDs, settings)
			s.Brokers[id].HideNotification()
			s.Brokers[id].SendToast(models.NewInfoToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "", ""))
		}(userID, settings)

		w.WriteHeader(http.StatusAccepted)
	}
}
