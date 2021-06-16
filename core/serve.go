package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/reaper47/recipe-hunter/api"
	"github.com/reaper47/recipe-hunter/config"
	"github.com/reaper47/recipe-hunter/repository"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

// Serve starts the web server at the address
// specified in the configration file.
func Serve() {
	env := InitEnv(repository.Db())

	interval := config.Config.IndexIntervalToDuration()
	log.Printf("Database indexing has been scheduled for every %v\n", interval)
	schedule(Index, interval)
	Index()

	r := createRouter(env)
	srv := createServer(r)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	wait := time.Second * time.Duration(config.Config.Wait)
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down recipya server...")
	os.Exit(0)
}

func createRouter(env *Env) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.NotFoundHandler = http.HandlerFunc(handle404)

	apiRootRouter := r.PathPrefix(api.ApiUrlSuffix).Subrouter()
	initRecipesRoutes(apiRootRouter, env)

	spaRouter := spaHandler{staticPath: config.Config.WebAppDir, indexPath: "index.html"}
	r.PathPrefix("/").Handler(spaRouter)

	r.Use(mux.CORSMethodMiddleware(r))
	return r
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func createServer(r *mux.Router) *http.Server {
	addr := ":" + strconv.Itoa(config.Config.Port)
	log.Println("Server started @ " + addr)
	return &http.Server{
		Addr:         addr,
		Handler:      r,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}
