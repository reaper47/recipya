package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/docs"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/jobs"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/services"
	_ "github.com/reaper47/recipya/internal/templates" // Need to initialize the templates package.
	"github.com/reaper47/recipya/web/static"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func init() {
	SessionData.Data = make(map[uuid.UUID]int64)
}

// NewServer creates a Server.
func NewServer(repo services.RepositoryService) *Server {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic("scraper not initialized")
	}

	srv := &Server{
		Brokers:      make(map[int64]*models.Broker),
		Email:        services.NewEmailService(),
		Files:        services.NewFilesService(),
		Integrations: services.NewIntegrationsService(&http.Client{}),
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		Repository:   repo,
		Scraper: scraper.NewScraper(&http.Client{
			Jar: jar,
		}),
	}
	srv.mountHandlers()

	f, err := os.Open(filepath.Join(filepath.Dir(app.DBBasePath), "sessions.csv"))
	if err == nil {
		slog.Info("Restoring user sessions")
		SessionData.Load(f)
		_ = f.Close()
		os.Remove(f.Name())
	}

	return srv
}

// Server is the web application's configuration object.
type Server struct {
	Brokers      map[int64]*models.Broker
	Email        services.EmailService
	Files        services.FilesService
	Integrations services.IntegrationsService
	Logger       *slog.Logger
	Repository   services.RepositoryService
	Router       *http.ServeMux
	Scraper      scraper.IScraper
}

func (s *Server) mountHandlers() {
	mux := http.NewServeMux()

	// Static files routes
	subFS, _ := fs.Sub(docs.FS, "website/public")
	mux.HandleFunc("GET /guide", http.StripPrefix("/guide", http.FileServerFS(subFS)).ServeHTTP)
	mux.HandleFunc("GET /guide/*", http.StripPrefix("/guide", http.FileServerFS(subFS)).ServeHTTP)
	mux.HandleFunc("GET /guide/login", guideLoginHandler)
	mux.HandleFunc("GET /static/*", http.StripPrefix("/static", http.FileServerFS(static.FS)).ServeHTTP)
	mux.HandleFunc("GET /data/images/*", http.StripPrefix("/data/images", http.FileServer(http.Dir(app.ImagesDir))).ServeHTTP)
	mux.HandleFunc("GET /*", notFoundHandler)

	// General routes
	mux.HandleFunc("GET /{$}", s.indexHandler)
	mux.Handle("GET /download/{tmpFile}", s.mustBeLoggedInMiddleware(s.downloadHandler()))
	mux.Handle("GET /user-initials", s.mustBeLoggedInMiddleware(s.userInitialsHandler()))
	mux.Handle("GET /update", s.mustBeLoggedInMiddleware(s.updateHandler()))
	mux.Handle("GET /ws", s.mustBeLoggedInMiddleware(s.wsHandler()))

	// Admin routes
	adminMiddleware := func(next http.Handler) http.Handler { return s.mustBeLoggedInMiddleware(s.onlyAdminMiddleware(next)) }
	mux.Handle("GET /admin", adminMiddleware(s.adminHandler()))
	mux.Handle("POST /admin/users", adminMiddleware(s.adminUsersPostHandler()))
	mux.Handle("DELETE /admin/users/{email}", adminMiddleware(s.adminUsersDeleteHandler()))

	// Auth routes
	withAuthRegister := func(next http.Handler) http.Handler {
		return s.redirectIfLoggedInMiddleware(redirectIfNoSignupsMiddleware(next))
	}
	withLog := func(next http.Handler) http.Handler {
		return s.mustBeLoggedInMiddleware(s.loggingMiddleware(next))
	}
	mux.Handle("POST /auth/change-password", s.mustBeLoggedInMiddleware(s.changePasswordHandler()))
	mux.HandleFunc("GET /auth/confirm", s.confirmHandler)
	mux.HandleFunc("GET /auth/forgot-password", s.forgotPasswordHandler)
	mux.HandleFunc("POST /auth/forgot-password", s.forgotPasswordPostHandler)
	mux.HandleFunc("GET /auth/forgot-password/reset", forgotPasswordResetHandler)
	mux.HandleFunc("POST /auth/forgot-password/reset", s.forgotPasswordResetPostHandler)
	mux.Handle("GET /auth/login", s.redirectIfLoggedInMiddleware(loginHandler()))
	mux.Handle("POST /auth/login", s.redirectIfLoggedInMiddleware(s.loginPostHandler()))
	mux.Handle("GET /auth/register", withAuthRegister(registerHandler()))
	mux.Handle("POST /auth/register", withAuthRegister(s.registerPostHandler()))
	mux.HandleFunc("POST /auth/logout", s.logoutHandler)
	mux.Handle("DELETE /auth/user", s.mustBeLoggedInMiddleware(s.deleteUserHandler()))

	// Cookbooks routes
	mux.Handle("GET /cookbooks", s.mustBeLoggedInMiddleware(s.cookbooksHandler()))
	mux.Handle("POST /cookbooks", withLog(s.cookbooksPostHandler()))
	mux.Handle("GET /cookbooks/{id}", s.mustBeLoggedInMiddleware(s.cookbooksGetCookbookHandler()))
	mux.Handle("POST /cookbooks/{id}", withLog(s.cookbookPostCookbookHandler()))
	mux.Handle("DELETE /cookbooks/{id}", withLog(s.cookbooksDeleteCookbookHandler()))
	mux.Handle("GET /cookbooks/{id}/download", s.mustBeLoggedInMiddleware(s.cookbooksDownloadCookbookHandler()))
	mux.Handle("PUT /cookbooks/{id}/image", withLog(s.cookbooksImagePostCookbookHandler()))
	mux.Handle("PUT /cookbooks/{id}/reorder", withLog(s.cookbooksPostCookbookReorderHandler()))
	mux.Handle("DELETE /cookbooks/{id}/recipes/{recipeID}", s.mustBeLoggedInMiddleware(s.cookbooksDeleteCookbookRecipeHandler()))
	mux.Handle("GET /cookbooks/{id}/recipes/search", s.mustBeLoggedInMiddleware(s.cookbooksRecipesSearchHandler()))
	mux.Handle("POST /cookbooks/{id}/share", withLog(s.cookbookSharePostHandler()))

	// Integrations routes
	mux.Handle("POST /integrations/import", withLog(s.integrationsImport()))

	// Recipes routes
	mux.Handle("GET /recipes", s.mustBeLoggedInMiddleware(s.recipesHandler()))
	mux.Handle("GET /recipes/{id}", s.mustBeLoggedInMiddleware(s.recipesViewHandler()))
	mux.Handle("DELETE /recipes/{id}", withLog(s.recipeDeleteHandler()))
	mux.Handle("GET /recipes/{id}/scale", s.mustBeLoggedInMiddleware(s.recipeScaleHandler()))
	mux.Handle("POST /recipes/{id}/share", withLog(s.recipeSharePostHandler()))
	mux.Handle("GET /recipes/{id}/edit", s.mustBeLoggedInMiddleware(s.recipesEditHandler()))
	mux.Handle("PUT /recipes/{id}/edit", withLog(s.recipesEditPostHandler()))
	mux.Handle("GET /recipes/add", s.mustBeLoggedInMiddleware(recipesAddHandler()))
	mux.Handle("POST /recipes/add/import", withLog(s.recipesAddImportHandler()))
	mux.Handle("GET /recipes/add/manual", s.mustBeLoggedInMiddleware(s.recipeAddManualHandler()))
	mux.Handle("POST /recipes/add/manual", withLog(s.recipeAddManualPostHandler()))
	mux.Handle("POST /recipes/add/manual/ingredient", s.mustBeLoggedInMiddleware(recipeAddManualIngredientHandler()))
	mux.Handle("POST /recipes/add/manual/ingredient/{entry}", s.mustBeLoggedInMiddleware(recipeAddManualIngredientDeleteHandler()))
	mux.Handle("POST /recipes/add/manual/instruction", s.mustBeLoggedInMiddleware(recipeAddManualInstructionHandler()))
	mux.Handle("POST /recipes/add/manual/instruction/{entry}", s.mustBeLoggedInMiddleware(recipeAddManualInstructionDeleteHandler()))
	mux.Handle("POST /recipes/add/ocr", withLog(s.recipesAddOCRHandler()))
	mux.Handle("POST /recipes/add/website", withLog(s.recipesAddWebsiteHandler()))
	mux.Handle("GET /recipes/search", s.mustBeLoggedInMiddleware(s.recipesSearchHandler()))
	mux.Handle("GET /recipes/supported-applications", s.mustBeLoggedInMiddleware(s.recipesSupportedApplicationsHandler()))
	mux.Handle("GET /recipes/supported-websites", s.mustBeLoggedInMiddleware(s.recipesSupportedWebsitesHandler()))

	// Reports routes
	mux.Handle("GET /reports", s.mustBeLoggedInMiddleware(s.reportsHandler()))
	mux.Handle("GET /reports/{id}", s.mustBeLoggedInMiddleware(s.reportsReportHandler()))

	// Settings routes
	mux.Handle("GET /settings", s.mustBeLoggedInMiddleware(s.settingsHandler()))
	mux.Handle("GET /settings/export/recipes", s.mustBeLoggedInMiddleware(s.settingsExportRecipesHandler()))
	mux.Handle("POST /settings/calculate-nutrition", withLog(s.settingsCalculateNutritionPostHandler()))
	mux.Handle("PUT /settings/config", withLog(s.onlyAdminMiddleware(settingsConfigPutHandler())))
	mux.Handle("POST /settings/convert-automatically", withLog(s.settingsConvertAutomaticallyPostHandler()))
	mux.Handle("POST /settings/measurement-system", withLog(s.settingsMeasurementSystemsPostHandler()))
	mux.Handle("POST /settings/backups/restore", withLog(s.settingsBackupsRestoreHandler()))
	mux.Handle("GET /settings/tabs/advanced", s.mustBeLoggedInMiddleware(s.settingsTabsAdvancedHandler()))
	mux.Handle("GET /settings/tabs/profile", s.mustBeLoggedInMiddleware(settingsTabsProfileHandler()))
	mux.Handle("GET /settings/tabs/recipes", s.mustBeLoggedInMiddleware(s.settingsTabsRecipesHandler()))

	// Share routes
	mux.HandleFunc("GET /r/{id}", s.recipeShareHandler)
	mux.HandleFunc("GET /c/{id}", s.cookbookShareHandler)

	s.Router = mux
}

// Run starts the web server.
func (s *Server) Run() {
	handler := slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   filepath.Join(app.LogsDir, "recipya.log"),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}, nil)
	logger := slog.New(handler)
	s.Logger = logger
	slog.SetDefault(logger)

	httpServer := &http.Server{
		Addr:              "0.0.0.0:" + strconv.Itoa(app.Config.Server.Port),
		ErrorLog:          slog.NewLogLogger(handler, slog.LevelError),
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
				fmt.Println("Forcing exit as graceful shutdown timed out")
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

	if app.Config.Server.IsAutologin {
		err := s.Repository.InitAutologin()
		if err != nil {
			fmt.Printf("Could not initialize the autologin feature: %q\n", err)
			os.Exit(1)
		}
	}

	jobs.ScheduleCronJobs(s.Repository, s.Files, s.Email)

	fmt.Println("Serving HTTP server at address", app.Config.Address())
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-serverCtx.Done()
}
