package core

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/reaper47/recipe-hunter/api"
	"github.com/reaper47/recipe-hunter/config"
	"github.com/reaper47/recipe-hunter/repository"
)

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

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.Config.WebAppDir)))
	return r
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func createServer(r *mux.Router) *http.Server {
	addr := ":" + strconv.Itoa(config.Config.Port)
	log.Println("Server started @ " + addr)
	return &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
}
