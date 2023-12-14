package server

import (
	"fmt"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
)

func (s *Server) integrationsImportNextcloud(w http.ResponseWriter, r *http.Request) {
	if s.Brokers == nil {
		w.Header().Set("HX-Trigger", makeToast("Connection lost. Please reload page.", warningToast))
		w.WriteHeader(http.StatusInternalServerError)
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

	go func(userID int64) {
		err := s.Brokers[userID].SendProgressStatus("Contacting server...", true, 0, -1)
		if err != nil {
			fmt.Println(err)
			return
		}

		var (
			recipes   *models.Recipes
			processed int
			result    = make(chan models.Result)
			errors    = make(chan error, 1)
		)

		go func() {
			defer close(result)
			recipes, err = s.Integrations.NextcloudImport(baseURL, username, password, s.Files, result)
			if err != nil {
				errors <- err
			}
		}()

		select {
		case err = <-errors:
			_ = s.Brokers[userID].SendToast("Failed to import Nextcloud recipes.", "bg-error-500")
			fmt.Println(err)
			return
		case <-result:
			for r := range result {
				processed++
				err = s.Brokers[userID].SendProgress("Fetching recipes...", processed, r.Total*2)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		count := 0
		skipped := 0
		numRecipes := len(*recipes)
		for i, recipe := range *recipes {
			_ = s.Brokers[userID].SendProgress("Adding to collection...", i+numRecipes, numRecipes*2)
			c := recipe.Copy()
			_, err = s.Repository.AddRecipe(&c, userID)
			if err != nil {
				skipped++
				continue
			}
			count++
		}

		_ = s.Brokers[userID].SendProgressStatus("Finished", false, 0, 100)
		_ = s.Brokers[userID].SendToast(fmt.Sprintf("Imported %d recipes. Skipped %d.", count, skipped), "bg-blue-500")
	}(getUserID(r))

	w.WriteHeader(http.StatusAccepted)
}
