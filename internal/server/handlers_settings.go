package server

import (
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

	go func() {
		broker := s.Brokers[userID]

		qType := r.URL.Query().Get("type")
		fileType := models.NewFileType(qType)
		if fileType == models.InvalidFileType {
			_ = broker.SendToast("Invalid export file format.", "bg-red-500")
			return
		}

		err := broker.SendProgressStatus("Preparing...", true, 0, -1)
		if err != nil {
			fmt.Println(err)
			return
		}

		recipes := s.Repository.RecipesAll(userID)
		if len(recipes) == 0 {
			_ = broker.SendToast("No recipes in database.", "bg-yellow-500")
			return
		}

		data, err := s.Files.ExportRecipes(recipes, fileType, s.Brokers[userID])
		if err != nil {
			fmt.Println(err)
			return
		}

		err = broker.SendProgressStatus("Finished", false, 0, 100)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = broker.SendFile("recipes_"+qType+".zip", data)
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
