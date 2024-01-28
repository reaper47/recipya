package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/docs"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/jobs"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services"
	_ "github.com/reaper47/recipya/internal/templates" // Need to initialize the templates package.
	"github.com/reaper47/recipya/static"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var imagesDir string

func init() {
	SessionData = make(map[uuid.UUID]int64)

	exe, err := os.Executable()
	if err != nil {
		return
	}
	imagesDir = filepath.Join(filepath.Dir(exe), "data", "images")
}

// NewServer creates a Server.
func NewServer(repository services.RepositoryService, email services.EmailService, files services.FilesService, integrations services.IntegrationsService) *Server {
	srv := &Server{
		Brokers:      make(map[int64]*models.Broker),
		Repository:   repository,
		Email:        email,
		Files:        files,
		Integrations: integrations,
	}
	srv.mountHandlers()
	return srv
}

// Server is the web application's configuration object.
type Server struct {
	Brokers      map[int64]*models.Broker
	Router       *chi.Mux
	Repository   services.RepositoryService
	Email        services.EmailService
	Files        services.FilesService
	Integrations services.IntegrationsService
}

func (s *Server) mountHandlers() {
	r := chi.NewRouter()

	r.Mount("/debug", middleware.Profiler())
	r.Get("/", s.indexHandler)

	subFS, _ := fs.Sub(docs.FS, "website/public")
	r.Get("/guide*", http.StripPrefix("/guide", http.FileServer(http.FS(subFS))).ServeHTTP)
	r.Get("/guide/login", guideLoginHandler)
	r.Get("/static/*", http.StripPrefix("/static", http.FileServer(http.FS(static.FS))).ServeHTTP)
	r.Get("/data/images/*", http.StripPrefix("/data/images", http.FileServer(http.Dir(imagesDir))).ServeHTTP)

	r.NotFound(notFoundHandler)

	r.Route("/r", func(r chi.Router) {
		r.Get("/{id:[1-9]([0-9])*}", s.recipesViewShareHandler)
		r.Get("/{id:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}", s.recipeShareHandler)
	})
	r.Get("/c/{id:^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}", s.cookbookShareHandler)

	r.Route("/auth", func(r chi.Router) {
		r.With(s.mustBeLoggedInMiddleware).Post("/change-password", s.changePasswordHandler)
		r.Get("/confirm", s.confirmHandler)

		r.Route("/forgot-password", func(r chi.Router) {
			r.Get("/", s.forgotPasswordHandler)
			r.Post("/", s.forgotPasswordPostHandler)

			r.Route("/reset", func(r chi.Router) {
				r.Get("/", forgotPasswordResetHandler)
				r.Post("/", s.forgotPasswordResetPostHandler)
			})
		})

		r.Route("/login", func(r chi.Router) {
			r.Use(s.redirectIfLoggedInMiddleware)

			r.Get("/", loginHandler)
			r.Post("/", s.loginPostHandler)
		})

		r.Route("/register", func(r chi.Router) {
			r.Use(s.redirectIfLoggedInMiddleware)

			r.Get("/", registerHandler)
			r.Post("/", s.registerPostHandler)
			r.Post("/validate-password", s.registerPostPasswordHandler)
		})

		r.Post("/logout", s.logoutHandler)

		r.Group(func(r chi.Router) {
			r.Use(s.mustBeLoggedInMiddleware)

			r.Delete("/user", s.deleteUserHandler)
		})
	})

	r.Route("/cookbooks", func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Get("/", s.cookbooksHandler)
		r.Post("/", s.cookbooksPostHandler)
		r.Post("/recipes/search", s.cookbooksRecipesSearchPostHandler)

		r.Route("/{id:[1-9]([0-9])*}", func(r chi.Router) {
			r.Get("/", s.cookbooksGetCookbookHandler)
			r.Delete("/", s.cookbooksDeleteCookbookHandler)
			r.Post("/", s.cookbookPostCookbookHandler)
			r.Get("/download", s.cookbooksDownloadCookbookHandler)
			r.Put("/image", s.cookbooksImagePostCookbookHandler)
			r.Put("/reorder", s.cookbooksPostCookbookReorderHandler)
			r.Delete("/recipes/{recipeID:[1-9]([0-9])*}", s.cookbooksDeleteCookbookRecipeHandler)
			r.Post("/share", s.cookbookSharePostHandler)
		})
	})

	r.Route("/integrations", func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Route("/import", func(r chi.Router) {
			r.Post("/nextcloud", s.integrationsImportNextcloud)
		})
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Get("/", s.recipesHandler)

		r.Route("/{id:[1-9]([0-9])*}", func(r chi.Router) {
			r.Get("/", s.recipesViewHandler)
			r.Delete("/", s.recipeDeleteHandler)
			r.Get("/scale", s.recipeScaleHandler)
			r.Post("/share", s.recipeSharePostHandler)

			r.Route("/edit", func(r chi.Router) {
				r.Get("/", s.recipesEditHandler)
				r.Put("/", s.recipesEditPostHandler)
			})
		})

		r.Route("/add", func(r chi.Router) {
			r.Get("/", recipesAddHandler)
			r.Post("/import", s.recipesAddImportHandler)

			r.Route("/manual", func(r chi.Router) {
				r.Get("/", s.recipeAddManualHandler)
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

			r.Post("/ocr", s.recipesAddOCRHandler)
			r.Post("/request-website", s.recipesAddRequestWebsiteHandler)
			r.Post("/website", s.recipesAddWebsiteHandler)
		})

		r.Post("/search", s.recipesSearchHandler)
		r.Get("/supported-websites", s.recipesSupportedWebsitesHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(s.mustBeLoggedInMiddleware)

		r.Route("/settings", func(r chi.Router) {
			r.Get("/", s.settingsHandler)

			r.Route("/export", func(r chi.Router) {
				r.Get("/recipes", s.settingsExportRecipesHandler)
			})

			r.Post("/calculate-nutrition", s.settingsCalculateNutritionPostHandler)
			r.Post("/convert-automatically", s.settingsConvertAutomaticallyPostHandler)
			r.Post("/measurement-system", s.settingsMeasurementSystemsPostHandler)

			r.Route("/backups", func(r chi.Router) {
				r.Post("/restore", s.settingsBackupsRestoreHandler)
			})

			r.Route("/tabs", func(r chi.Router) {
				r.Get("/advanced", s.settingsTabsAdvancedHandler)
				r.Get("/profile", settingsTabsProfileHandler)
				r.Get("/recipes", s.settingsTabsRecipesHandler)
			})
		})

		r.Get("/download/{tmpFile}", s.downloadHandler)
		r.Get("/user-initials", s.userInitialsHandler)
		r.Get("/ws", s.wsHandler)
	})

	s.Router = r
}

// Run starts the web server.
func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:              "0.0.0.0:" + strconv.Itoa(app.Config.Server.Port),
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

		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		serverStopCtx()
	}()

	err := s.Repository.InitAutologin()
	if err != nil {
		fmt.Printf("Could not initialize the autologin feature: %q\n", err)
		os.Exit(1)
	}

	jobs.ScheduleCronJobs(s.Repository, s.Files, s.Email)

	fmt.Printf("Serving on %s\n", app.Config.Address())
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-serverCtx.Done()
}
