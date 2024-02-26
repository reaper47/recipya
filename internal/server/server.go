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
func NewServer(
	repository services.RepositoryService, email services.EmailService,
	files services.FilesService, integrations services.IntegrationsService,
	scraper scraper.IScraper,
) *Server {
	srv := &Server{
		Brokers:      make(map[int64]*models.Broker),
		Email:        email,
		Files:        files,
		Integrations: integrations,
		Repository:   repository,
		Scraper:      scraper,
	}
	srv.mountHandlers()
	return srv
}

// Server is the web application's configuration object.
type Server struct {
	Brokers      map[int64]*models.Broker
	Email        services.EmailService
	Files        services.FilesService
	Integrations services.IntegrationsService
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
	mux.HandleFunc("GET /data/images/*", http.StripPrefix("/data/images", http.FileServer(http.Dir(imagesDir))).ServeHTTP)
	mux.HandleFunc("GET /*", notFoundHandler)

	// General routes
	mux.HandleFunc("GET /{$}", s.indexHandler)
	mux.Handle("GET /download/{tmpFile}", s.mustBeLoggedInMiddleware(s.downloadHandler()))
	mux.Handle("GET /user-initials", s.mustBeLoggedInMiddleware(s.userInitialsHandler()))
	mux.Handle("GET /ws", s.mustBeLoggedInMiddleware(s.wsHandler()))

	// Admin routes
	adminMiddleware := func(next http.Handler) http.Handler { return s.mustBeLoggedInMiddleware(s.onlyAdminMiddleware(next)) }
	mux.Handle("GET /admin", adminMiddleware(s.adminHandler()))
	mux.Handle("POST /admin/users", adminMiddleware(s.adminUsersPostHandler()))
	mux.Handle("DELETE /admin/users/{email}", adminMiddleware(s.adminUsersDeleteHandler()))

	// Auth routes
	authRegisterMiddleware := func(next http.Handler) http.Handler {
		return s.redirectIfLoggedInMiddleware(redirectIfNoSignupsMiddleware(next))
	}
	mux.Handle("POST /auth/change-password", s.mustBeLoggedInMiddleware(s.changePasswordHandler()))
	mux.HandleFunc("GET /auth/confirm", s.confirmHandler)
	mux.HandleFunc("GET /auth/forgot-password", s.forgotPasswordHandler)
	mux.HandleFunc("POST /auth/forgot-password", s.forgotPasswordPostHandler)
	mux.HandleFunc("GET /auth/forgot-password/reset", forgotPasswordResetHandler)
	mux.HandleFunc("POST /auth/forgot-password/reset", s.forgotPasswordResetPostHandler)
	mux.Handle("GET /auth/login", s.redirectIfLoggedInMiddleware(loginHandler()))
	mux.Handle("POST /auth/login", s.redirectIfLoggedInMiddleware(s.loginPostHandler()))
	mux.Handle("GET /auth/register", authRegisterMiddleware(registerHandler()))
	mux.Handle("POST /auth/register", authRegisterMiddleware(s.registerPostHandler()))
	mux.HandleFunc("POST /auth/logout", s.logoutHandler)
	mux.Handle("DELETE /auth/user", s.mustBeLoggedInMiddleware(s.deleteUserHandler()))

	// Cookbooks routes
	mux.Handle("GET /cookbooks", s.mustBeLoggedInMiddleware(s.cookbooksHandler()))
	mux.Handle("POST /cookbooks", s.mustBeLoggedInMiddleware(s.cookbooksPostHandler()))
	mux.Handle("POST /cookbooks/recipes/search", s.mustBeLoggedInMiddleware(s.cookbooksRecipesSearchPostHandler()))
	mux.Handle("GET /cookbooks/{id}", s.mustBeLoggedInMiddleware(s.cookbooksGetCookbookHandler()))
	mux.Handle("POST /cookbooks/{id}", s.mustBeLoggedInMiddleware(s.cookbookPostCookbookHandler()))
	mux.Handle("DELETE /cookbooks/{id}", s.mustBeLoggedInMiddleware(s.cookbooksDeleteCookbookHandler()))
	mux.Handle("GET /cookbooks/{id}/download", s.mustBeLoggedInMiddleware(s.cookbooksDownloadCookbookHandler()))
	mux.Handle("PUT /cookbooks/{id}/image", s.mustBeLoggedInMiddleware(s.cookbooksImagePostCookbookHandler()))
	mux.Handle("PUT /cookbooks/{id}/reorder", s.mustBeLoggedInMiddleware(s.cookbooksPostCookbookReorderHandler()))
	mux.Handle("DELETE /cookbooks/{id}/recipes/{recipeID}", s.mustBeLoggedInMiddleware(s.cookbooksDeleteCookbookRecipeHandler()))
	mux.Handle("POST /cookbooks/{id}/share", s.mustBeLoggedInMiddleware(s.cookbookSharePostHandler()))

	// Integrations routes
	mux.Handle("POST /integrations/import/nextcloud", s.mustBeLoggedInMiddleware(s.integrationsImportNextcloud()))

	// Recipes routes
	mux.Handle("GET /recipes", s.mustBeLoggedInMiddleware(s.recipesHandler()))
	mux.Handle("GET /recipes/{id}", s.mustBeLoggedInMiddleware(s.recipesViewHandler()))
	mux.Handle("DELETE /recipes/{id}", s.mustBeLoggedInMiddleware(s.recipeDeleteHandler()))
	mux.Handle("GET /recipes/{id}/scale", s.mustBeLoggedInMiddleware(s.recipeScaleHandler()))
	mux.Handle("POST /recipes/{id}/share", s.mustBeLoggedInMiddleware(s.recipeSharePostHandler()))
	mux.Handle("GET /recipes/{id}/edit", s.mustBeLoggedInMiddleware(s.recipesEditHandler()))
	mux.Handle("PUT /recipes/{id}/edit", s.mustBeLoggedInMiddleware(s.recipesEditPostHandler()))
	mux.Handle("GET /recipes/add", s.mustBeLoggedInMiddleware(recipesAddHandler()))
	mux.Handle("POST /recipes/add/import", s.mustBeLoggedInMiddleware(s.recipesAddImportHandler()))
	mux.Handle("GET /recipes/add/manual", s.mustBeLoggedInMiddleware(s.recipeAddManualHandler()))
	mux.Handle("POST /recipes/add/manual", s.mustBeLoggedInMiddleware(s.recipeAddManualPostHandler()))
	mux.Handle("POST /recipes/add/manual/ingredient", s.mustBeLoggedInMiddleware(recipeAddManualIngredientHandler()))
	mux.Handle("POST /recipes/add/manual/ingredient/{entry}", s.mustBeLoggedInMiddleware(recipeAddManualIngredientDeleteHandler()))
	mux.Handle("POST /recipes/add/manual/instruction", s.mustBeLoggedInMiddleware(recipeAddManualInstructionHandler()))
	mux.Handle("POST /recipes/add/manual/instruction/{entry}", s.mustBeLoggedInMiddleware(recipeAddManualInstructionDeleteHandler()))
	mux.Handle("POST /recipes/add/ocr", s.mustBeLoggedInMiddleware(s.recipesAddOCRHandler()))
	mux.Handle("POST /recipes/add/request-website", s.mustBeLoggedInMiddleware(s.recipesAddRequestWebsiteHandler()))
	mux.Handle("POST /recipes/add/website", s.mustBeLoggedInMiddleware(s.recipesAddWebsiteHandler()))
	mux.Handle("GET /recipes/search", s.mustBeLoggedInMiddleware(s.recipesSearchHandler()))
	mux.Handle("POST /recipes/search", s.mustBeLoggedInMiddleware(s.recipesSearchPostHandler()))
	mux.Handle("GET /recipes/supported-websites", s.mustBeLoggedInMiddleware(s.recipesSupportedWebsitesHandler()))

	// Reports routes
	mux.Handle("GET /reports", s.mustBeLoggedInMiddleware(s.reportsHandler()))
	mux.Handle("GET /reports/{id}", s.mustBeLoggedInMiddleware(s.reportsReportHandler()))

	// Settings routes
	mux.Handle("GET /settings", s.mustBeLoggedInMiddleware(s.settingsHandler()))
	mux.Handle("GET /settings/export/recipes", s.mustBeLoggedInMiddleware(s.settingsExportRecipesHandler()))
	mux.Handle("POST /settings/calculate-nutrition", s.mustBeLoggedInMiddleware(s.settingsCalculateNutritionPostHandler()))
	mux.Handle("POST /settings/convert-automatically", s.mustBeLoggedInMiddleware(s.settingsConvertAutomaticallyPostHandler()))
	mux.Handle("POST /settings/measurement-system", s.mustBeLoggedInMiddleware(s.settingsMeasurementSystemsPostHandler()))
	mux.Handle("POST /settings/backups/restore", s.mustBeLoggedInMiddleware(s.settingsBackupsRestoreHandler()))
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

	if app.Config.Server.IsAutologin {
		err := s.Repository.InitAutologin()
		if err != nil {
			fmt.Printf("Could not initialize the autologin feature: %q\n", err)
			os.Exit(1)
		}
	}

	jobs.ScheduleCronJobs(s.Repository, s.Files, s.Email)

	fmt.Printf("Serving on %s\n", app.Config.Address())
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	<-serverCtx.Done()
}
