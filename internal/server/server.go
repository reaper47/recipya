package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/services"
	_ "github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/static"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func init() {
	SessionData = make(map[uuid.UUID]int64)
}

// NewServer creates a Server.
func NewServer(repository services.RepositoryService, email services.EmailService, files services.FilesService) *Server {
	srv := &Server{
		Repository: repository,
		Email:      email,
		Files:      files,
	}
	srv.mountHandlers()
	return srv
}

// Server is the web application's configuration object.
type Server struct {
	Router     *chi.Mux
	Repository services.RepositoryService
	Email      services.EmailService
	Files      services.FilesService
}

func (s *Server) mountHandlers() {
	r := chi.NewRouter()

	r.Get("/", s.indexHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/confirm", s.confirmHandler)

		r.Route("/login", func(r chi.Router) {
			r.Use(s.redirectIfLoggedInMiddleware)

			r.Get("/", loginHandler)
			r.Post("/", s.loginPostHandler)
		})

		r.Route("/register", func(r chi.Router) {
			r.Use(s.redirectIfLoggedInMiddleware)

			r.Get("/", registerHandler)
			r.Post("/", s.registerPostHandler)
			r.Post("/validate-email", s.registerPostEmailHandler)
			r.Post("/validate-password", s.registerPostPasswordHandler)
		})

		r.Post("/logout", s.logoutHandler)
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Route("/add", func(r chi.Router) {
			r.Get("/", recipesAddHandler)
			r.Post("/import", s.recipesAddImportHandler)

			r.Route("/manual", func(r chi.Router) {
				r.Get("/", recipeAddManualHandler)
				r.Post("/", s.recipeAddManualPostHandler)

				r.Route("/ingredient", func(r chi.Router) {
					r.Post("/", recipeAddManualIngredientHandler)
					r.Post("/{entry:[1-9]([0-9])*}", recipeAddManualIngredientDeleteHandler)
				})

				r.Route("/instruction", func(r chi.Router) {
					r.Post("/", recipeAddManualInstructionHandler)
					r.Post("/{entry:[1-9]([0-9])*}", recipeAddManualInstructionDeleteHandler)
				})
			})

			r.Post("/request-website", s.recipesAddRequestWebsiteHandler)
			r.Post("/website", s.recipesAddWebsiteHandler)
		})

		r.Route("/supported-websites", func(r chi.Router) {
			r.Get("/", s.recipesSupportedWebsitesHandler)
			r.Post("/", s.recipesSupportedWebsitesPostHandler)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Get("/avatar-dropdown", avatarDropdownHandler)
		r.Get("/settings", s.settingsHandler)
		r.Get("/user-initials", s.userInitialsHandler)
	})

	r.NotFound(notFoundHandler)

	staticFS := http.FileServer(http.FS(static.FS))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", staticFS).ServeHTTP(w, r)
	})

	s.Router = r
}

// Run starts the web server.
func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:              "0.0.0.0:" + strconv.Itoa(app.Config.Port),
		Handler:           s.Router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       1 * time.Minute,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				fmt.Println("forcing exit as graceful shutdown timed out")
				os.Exit(1)
			}
		}()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		serverStopCtx()
	}()

	fmt.Printf("Serving on %s\n", app.Config.Address())
	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-serverCtx.Done()
}
