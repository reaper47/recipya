package server

import (
	"errors"
	"github.com/reaper47/recipya/web"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"syscall"
	"time"

	"github.com/coder/websocket"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/web/components"
)

func (s *Server) downloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := r.PathValue("tmpFile")
		data, err := s.Files.ReadTempFile(file)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(data))
		w.Header().Set("Content-Disposition", `attachment; filename="`+file+`"`)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		_, _ = w.Write(data)
	}
}

func (s *Server) fetchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")
		userID := getUserID(r)

		parsed, err := url.Parse(rawURL)
		if err != nil || rawURL == "" || !slices.Contains([]string{"http", "https"}, parsed.Scheme) {
			s.Brokers.SendToast(models.NewErrorGeneralToast("Invalid URL."), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := http.Get(rawURL)
		if err != nil {
			s.Brokers.SendToast(models.NewErrorGeneralToast("Could not fetch URL."), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer res.Body.Close()

		w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
		io.Copy(w, res.Body)
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if app.Config.Server.IsAutologin || isAuthenticated(r, s.Repository.GetAuthToken) {
		middleware := s.mustBeLoggedInMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.recipesHandler().ServeHTTP(w, r)
		}))
		middleware.ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = components.SimplePage("Page Not Found", "The page you requested to view is not found. Please go back to the main page.").Render(r.Context(), w)
}

func (s *Server) placeholderPostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		name := r.FormValue("name")
		if name != "recipe" && name != "cookbook" {
			s.Brokers.SendToast(models.NewErrorFormToast("Only the recipe or cookbook placeholder can be updated."), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<24)

		err := r.ParseMultipartForm(1 << 24)
		if err != nil {
			msg := "Could not parse the form."
			slog.Error(msg, userIDAttr, "error", err)
			s.Brokers.SendToast(models.NewErrorFormToast(msg), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		images, ok := r.MultipartForm.File["images"]
		if !ok {
			msg := "Could not retrieve the image from the form."
			slog.Error(msg, userIDAttr, "error", err)
			s.Brokers.SendToast(models.NewErrorFormToast(msg), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(images) == 0 {
			msg := "No image has been uploaded."
			slog.Error(msg, userIDAttr, "error", err)
			s.Brokers.SendToast(models.NewErrorFormToast(msg), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f, err := images[0].Open()
		if err != nil {
			msg := "Could not open the image from the form."
			slog.Error(msg, "error", err, userIDAttr)
			s.Brokers.SendToast(models.NewErrorFormToast(msg), userID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()

		imageUUID, err := s.Files.UploadImage(f)
		if err != nil {
			msg := "Error uploading image."
			slog.Error(msg, "error", err, userIDAttr)
			s.Brokers.SendToast(models.NewErrorFilesToast(msg), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		imageUUIDAttr := slog.String("imageUUID", imageUUID.String())

		from := filepath.Join(app.ImagesDir, imageUUID.String()+app.ImageExt)
		to := filepath.Join(app.ImagesDir, "Placeholders", "placeholder."+name+app.ImageExt)
		err = os.Rename(from, to)
		if err != nil {
			msg := "Error moving compressed placeholder image."
			slog.Error(msg, "error", err, imageUUIDAttr, userIDAttr)
			s.Brokers.SendToast(models.NewErrorFilesToast(msg), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Server) restorePlaceholderPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := "recipe"
		path := "static/img/recipes/placeholder.webp"
		if r.FormValue("name") == "cookbook" {
			name = "cookbook"
			path = "static/img/cookbooks-new/placeholder.webp"
		}

		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		openFile, err := web.StaticFS.Open(path)
		if err != nil {
			msg := "Error opening the public file."
			s.Brokers.SendToast(models.NewErrorGeneralToast(msg), userID)
			slog.Error(msg, "error", err, userIDAttr)
			return
		}
		defer openFile.Close()

		f, err := os.Create(filepath.Join(app.ImagesDir, "Placeholders", "placeholder."+name+".webp"))
		if err != nil {
			msg := "Error creating placeholder file."
			s.Brokers.SendToast(models.NewErrorGeneralToast(msg), userID)
			slog.Error(msg, "error", err, userIDAttr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		io.Copy(f, openFile)
	}
}

func (s *Server) userInitialsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(s.Repository.UserInitials(getUserID(r))))
	}
}

func (s *Server) updateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		err := s.Files.UpdateApp(app.Info.Version)
		if errors.Is(err, app.ErrNoUpdate) {
			s.Brokers.SendToast(models.NewWarningToast("", "No update available.", ""), userID)
			w.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
			msg := "Failed to update."
			slog.Error(msg, "error", err)
			s.Brokers.SendToast(models.NewErrorGeneralToast(msg), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.Brokers.SendToast(models.NewInfoToast("Software updated", "Application will reload in 5 seconds.", ""), userID)
		w.WriteHeader(http.StatusNoContent)

		go func() {
			slog.Info("Application will restart and data backed up.")

			err = s.Files.BackupGlobal()
			if err != nil {
				slog.Error("Backing up global data", "error", err)
				return
			}

			f, err := os.Create("sessions.csv")
			if err != nil {
				slog.Error("Failed to create file", "error", err)
				os.Exit(1)
			}
			defer f.Close()
			SessionData.Save(f)

			exe, err := os.Executable()
			if err != nil {
				slog.Error("Failed get executable path", "error", err)
				return
			}
			dir := filepath.Dir(exe)

			if runtime.GOOS == "windows" {
				err = exec.Command(filepath.Join(dir, "updater.exe")).Start()
				if err != nil {
					slog.Error("Failed to start application", "error", err)
					return
				}

				slog.Info("Started updater.exe. As you are on Windows, the running program can be found under Task Manager -> Details -> recipya.exe")
			} else {
				err = syscall.Exec(filepath.Join(dir, "recipya"), os.Args, os.Environ())
				if err != nil {
					slog.Error("Failed to start application", "error", err)
					return
				}
			}

			time.Sleep(250 * time.Millisecond)
			os.Exit(0)
		}()
	}
}

func (s *Server) updateCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)

		_, err := s.Repository.CheckUpdate(s.Files)
		if errors.Is(err, app.ErrNoUpdate) {
			s.Brokers.SendToast(models.NewWarningToast("", "No update available.", ""), userID)
		} else if err != nil {
			msg := "Failed to check update."
			slog.Error(msg, "error", err)
			s.Brokers.SendToast(models.NewErrorGeneralToast(msg), userID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		about := templates.NewAboutData()
		about.IsCheckUpdate = true

		_ = components.SettingsAbout(templates.Data{About: about}).Render(r.Context(), w)
	}
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewWarningToast("", "Could not upgrade connection.", "").Render())
			return
		}

		s.Brokers.Add(getUserID(r), c)
	}
}
