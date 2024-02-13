package server

import (
	"errors"
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
)

func (s *Server) integrationsImportNextcloud() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		_, found := s.Brokers[userID]
		if !found {
			w.Header().Set("HX-Trigger", makeToast("Connection lost. Please reload page.", warningToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		baseURL := r.FormValue("url")
		if username == "" || password == "" || baseURL == "" {
			w.Header().Set("HX-Trigger", makeToast("Invalid username, password or URL.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			fmt.Println(err)
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast("Failed to get user settings.", "bg-error-500")
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
				recipes, err = s.Integrations.NextcloudImport(baseURL, username, password, s.Files, progress)
				if err != nil {
					errs <- err
				}

				if recipes == nil {
					errs <- errors.New("recipes from Nextcloud is nil")
				}
			}()

			select {
			case err = <-errs:
				fmt.Println(err)
				s.Brokers[id].HideNotification()
				s.Brokers[id].SendToast("Failed to import Nextcloud recipes.", "bg-error-500")
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
					skipped++
					continue
				}
				recipeIDs = append(recipeIDs, recipeID)
				count++
			}

			s.Repository.CalculateNutrition(userID, recipeIDs, settings)
			s.Brokers[id].HideNotification()
			s.Brokers[id].SendToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "bg-blue-500")
		}(userID, settings)

		w.WriteHeader(http.StatusAccepted)
	}
}
