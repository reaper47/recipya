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

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			msg := "Failed to get user settings."
			slog.Error(msg, userIDAttr, "error", err)
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast(models.NewErrorDBToast(msg))
			return
		}

		integration := r.FormValue("integration")
		rawURL := r.FormValue("url")
		username := r.FormValue("username")
		password := r.FormValue("password")

		go func(id int64, us models.UserSettings) {
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
				}
			}

			recipes = slices.DeleteFunc(recipes, func(r models.Recipe) bool {
				return len(r.Instructions) == 0 && len(r.Ingredients) == 0 && r.Description == ""
			})

			count := 0
			skipped := 0
			numRecipes := len(recipes)
			recipeIDs := make([]int64, 0, numRecipes)

			for i, recipe := range recipes {
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

			slog.Info("Imported recipes", "integration", integration, userIDAttr, "count", count, "skipped", skipped)
			s.Repository.CalculateNutrition(userID, recipeIDs, settings)
			s.Brokers[id].HideNotification()
			s.Brokers[id].SendToast(models.NewInfoToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "", ""))
		}(userID, settings)

		w.WriteHeader(http.StatusAccepted)
	}
}
