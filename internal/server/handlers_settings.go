package server

import (
	"bytes"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"net/http"
	"slices"
)

func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "core", "settings", nil)
	} else {
		page := templates.SettingsPage
		templates.Render(w, page, templates.Data{
			About: templates.AboutData{
				Version: app.Version,
			},
			IsAuthenticated: true,
			Title:           page.Title(),
		})
	}
}

func (s *Server) settingsCalculateNutritionPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	isConvert := r.FormValue("calculate-nutrition") == "on"
	err := s.Repository.UpdateCalculateNutrition(userID, isConvert)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to set setting.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) settingsConvertAutomaticallyPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	isConvert := r.FormValue("convert") == "on"
	err := s.Repository.UpdateConvertMeasurementSystem(userID, isConvert)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to set setting.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) settingsExportRecipesHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	_, found := s.Brokers[userID]
	if !found {
		w.Header().Set("HX-Trigger", makeToast("Connection lost. Please reload page.", warningToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	if query == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qType := query.Get("type")
	fileType := models.NewFileType(qType)
	if fileType == models.InvalidFileType {
		w.Header().Set("HX-Trigger", makeToast("Invalid export file format.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		s.Brokers[userID].SendProgressStatus("Preparing...", true, 0, -1)

		recipes := s.Repository.RecipesAll(userID)
		if len(recipes) == 0 {
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast("No recipes in database.", "bg-orange-500")
			return
		}

		var (
			iter       = make(chan int)
			errors     = make(chan error, 1)
			numRecipes = len(recipes)
			err        error
		)

		var data *bytes.Buffer
		go func() {
			defer close(iter)
			data, err = s.Files.ExportRecipes(recipes, fileType, iter)
			if err != nil {
				errors <- err
				return
			}
		}()

		select {
		case err := <-errors:
			fmt.Println(err)
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast("Failed to export recipes.", "bg-error-500")
			return
		case <-iter:
			for value := range iter {
				s.Brokers[userID].SendProgress("Exporting recipes...", value, numRecipes)
			}
		}

		s.Brokers[userID].HideNotification()
		s.Brokers[userID].SendFile("recipes_"+qType+".zip", data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) settingsMeasurementSystemsPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	systems, settings, err := s.Repository.MeasurementSystems(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error fetching units systems.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	system := units.System(r.FormValue("system"))
	if !slices.Contains(systems, system) {
		w.Header().Set("HX-Trigger", makeToast("Measurement system does not exist.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if settings.MeasurementSystem == system {
		w.Header().Set("HX-Trigger", makeToast("System already set to "+system.String()+".", warningToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Repository.SwitchMeasurementSystem(system, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error switching units system.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", makeToast("Measurement system set to "+system.String()+".", infoToast))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) settingsBackupRestoreHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) settingsTabsAdvancedHandler(w http.ResponseWriter, r *http.Request) {
	backups := s.Files.Backups(getUserID(r))
	dates := make([]string, 0, len(backups))
	for _, backup := range backups {
		dates = append(dates, backup.Format("02 Jan 2006"))
	}
	templates.RenderComponent(w, "core", "settings-tabs-profile", dates)
}

func settingsTabsProfileHandler(w http.ResponseWriter, _ *http.Request) {
	templates.RenderComponent(w, "core", "settings-tabs-profile", nil)
}

func (s *Server) settingsTabsRecipesHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	systems, settings, err := s.Repository.MeasurementSystems(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error fetching units systems.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.RenderComponent(w, "core", "settings-tabs-recipes", templates.Data{
		Settings: templates.SettingsData{
			MeasurementSystems: systems,
			UserSettings:       settings,
		},
	})
}
