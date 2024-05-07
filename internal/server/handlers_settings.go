package server

import (
	"bytes"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/web/components"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"
)

func (s *Server) settingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data        templates.SettingsData
			isHxRequest = r.Header.Get("Hx-Request") == "true"
			userID      = getUserID(r)
		)

		if app.Config.Server.IsAutologin {
			systems, settings, err := s.Repository.MeasurementSystems(userID)
			if err != nil {
				msg := "Error fetching unit systems: " + err.Error()
				w.WriteHeader(http.StatusInternalServerError)
				if isHxRequest {
					s.Brokers.SendToast(models.NewErrorDBToast(msg), userID)
				} else {
					w.Write([]byte(msg))
				}
				return
			}
			data.UserSettings = settings
			data.MeasurementSystems = systems
		}

		_ = components.Settings(templates.Data{
			About:           templates.NewAboutData(),
			IsAdmin:         userID == 1,
			IsAutologin:     app.Config.Server.IsAutologin,
			IsAuthenticated: true,
			IsHxRequest:     isHxRequest,
			Settings:        data,
		}).Render(r.Context(), w)
	}
}

func (s *Server) settingsCalculateNutritionPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isConvert := r.FormValue("calculate-nutrition") == "on"
		err := s.Repository.UpdateCalculateNutrition(getUserID(r), isConvert)
		if err != nil {
			msg := "Failed to set setting."
			slog.Error(msg, "userID", getUserID(r), "error", err)
			s.Brokers.SendToast(models.NewErrorDBToast(msg), getUserID(r))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func settingsConfigPutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.Config.Update(app.ConfigFile{
			Email: app.ConfigEmail{
				From:           r.FormValue("email.from"),
				SendGridAPIKey: r.FormValue("email.apikey"),
			},
			Integrations: app.ConfigIntegrations{
				AzureDI: app.AzureDI{
					Key:      r.FormValue("integrations.ocr.key"),
					Endpoint: r.FormValue("integrations.ocr.url"),
				},
			},
			Server: app.ConfigServer{
				IsAutologin:  r.FormValue("server.autologin") == "on",
				IsNoSignups:  r.FormValue("server.noSignups") == "on",
				IsProduction: r.FormValue("server.production") == "on",
			},
		})
		if err != nil {
			msg := "Failed to update configuration."
			slog.Error(msg, "userID", getUserID(r), "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", models.NewInfoToast("Configuration updated.", "", "").Render())
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) settingsConvertAutomaticallyPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isConvert := r.FormValue("convert") == "on"
		err := s.Repository.UpdateConvertMeasurementSystem(getUserID(r), isConvert)
		if err != nil {
			msg := "Failed to set setting."
			slog.Error(msg, "userID", getUserID(r), "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) settingsExportRecipesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		if !s.Brokers.Has(userID) {
			w.Header().Set("HX-Trigger", models.NewWarningWSToast("Connection lost. Please reload page.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := r.URL.Query()
		if query == nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Could not parse query.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		qType := query.Get("type")
		fileType := models.NewFileType(qType)
		if fileType == models.InvalidFileType {
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast("Invalid export file format.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		go func() {
			s.Brokers.SendProgressStatus("Preparing...", true, 0, -1, userID)

			recipes := s.Repository.RecipesAll(userID)
			if len(recipes) == 0 {
				s.Brokers.HideNotification(userID)
				s.Brokers.SendToast(models.NewWarningToast("No recipes in database.", "", ""), userID)
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
				s.Brokers.HideNotification(userID)
				s.Brokers.SendToast(models.NewErrorFilesToast("Failed to export recipes."), userID)
				return
			case <-iter:
				for value := range iter {
					s.Brokers.SendProgress("Exporting recipes...", value, numRecipes, userID)
				}
			}

			s.Brokers.HideNotification(userID)
			s.Brokers.SendFile("recipes_"+qType+".zip", data, userID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()

		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) settingsMeasurementSystemsPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		systems, settings, err := s.Repository.MeasurementSystems(userID)
		if err != nil {
			msg := "Error fetching units systems."
			slog.Error(msg, "userID", userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		system := units.System(r.FormValue("system"))
		if !slices.Contains(systems, system) {
			msg := "Measurement system does not exist."
			slog.Error(msg, userIDAttr, "systems", systems, "settings", settings, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if settings.MeasurementSystem == system {
			msg := "System already set to " + system.String() + "."
			slog.Warn(msg, userIDAttr, "currentSystem", settings.MeasurementSystem, "system", system, "systems", systems)
			w.Header().Set("HX-Trigger", models.NewWarningToast(msg, "", "").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.Repository.SwitchMeasurementSystem(system, userID)
		if err != nil {
			msg := "Error switching units system."
			slog.Error(msg, userIDAttr, "system", system, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Switched measurement system", "from", settings.MeasurementSystem, "to", system)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) settingsBackupsRestoreHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		dateStr := r.FormValue("date")
		_, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			msg := dateStr + " is an invalid backup."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.Files.BackupUserData(s.Repository, userID)
		if err != nil {
			msg := "Failed to backup current data."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		backup, err := s.Files.ExtractUserBackup(dateStr, userID)
		if err != nil {
			msg := "Failed to extract backup."
			slog.Error(msg, userIDAttr, "date", dateStr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer os.RemoveAll(filepath.Dir(backup.ImagesPath))

		err = s.Repository.RestoreUserBackup(backup)
		if err != nil {
			msg := "Failed to restore backup."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := "Backup restored successfully."
		slog.Info(msg, userIDAttr, "date", dateStr)
		w.Header().Set("HX-Trigger", models.NewInfoToast(msg, "", "").Render())
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) settingsTabsAdvancedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backups := s.Files.Backups(getUserID(r))
		dates := make([]templates.Backup, 0, len(backups))
		for _, backup := range backups {
			dates = append(dates, templates.Backup{
				Display: backup.Format("02 Jan 2006"),
				Value:   backup.Format(time.DateOnly),
			})
		}

		c := app.Config
		if app.Config.Server.IsDemo {
			c.Email.From = "demo@demo.com"
			c.Email.SendGridAPIKey = "demo"
			c.Integrations.AzureDI.Key = "demo"
			c.Integrations.AzureDI.Endpoint = "https://www.example.com"
		}

		_ = components.SettingsTabsAdvanced(templates.SettingsData{
			Backups: dates,
			Config:  c,
		}, getUserID(r) == 1).Render(r.Context(), w)
	}
}

func settingsTabsProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = components.SettingsTabsProfile().Render(r.Context(), w)
	}
}

func (s *Server) settingsTabsRecipesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systems, settings, err := s.Repository.MeasurementSystems(getUserID(r))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error fetching units systems.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = components.SettingsTabsRecipes(templates.SettingsData{
			MeasurementSystems: systems,
			UserSettings:       settings,
		}).Render(r.Context(), w)
	}
}
