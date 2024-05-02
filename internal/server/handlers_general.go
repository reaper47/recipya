package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/web/components"
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
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

func fetchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawURL := r.URL.Query().Get("url")

		parsed, err := url.Parse(rawURL)
		if err != nil || rawURL == "" || !slices.Contains([]string{"http", "https"}, parsed.Scheme) {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Invalid URL.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := http.Get(rawURL)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Could not fetch URL.").Render())
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

	http.Redirect(w, r, "/guide", http.StatusSeeOther)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_ = components.SimplePage("Page Not Found", "The page you requested to view is not found. Please go back to the main page.").Render(r.Context(), w)
}

func (s *Server) userInitialsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(s.Repository.UserInitials(getUserID(r))))
	}
}

func (s *Server) updateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		err := s.Files.UpdateApp(app.Info.Version)
		if errors.Is(err, app.ErrNoUpdate) {
			w.Header().Set("HX-Trigger", models.NewWarningToast("", "No update available.", "").Render())
			w.WriteHeader(http.StatusNoContent)
			return
		} else if err != nil {
			msg := "Failed to update."
			slog.Error(msg, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Trigger", models.NewInfoToast("Software updated", "Application will reload in 5 seconds.", "").Render())
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

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewWarningToast("", "Could not upgrade connection.", "").Render())
			return
		}

		userID := getUserID(r)
		broker := models.NewBroker(userID, s.Brokers, ws)
		s.Brokers[userID] = broker
	}
}
