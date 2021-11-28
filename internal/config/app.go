package config

import (
	"log"
	"sync"

	"github.com/reaper47/recipya/internal/repository"
	db "github.com/reaper47/recipya/internal/repository/dbrepo"
)

// DBName is name the database.
const DBName = "recipya"

var (
	app  AppConfig
	once sync.Once
)

// AppConfig holds configuration data for the application.
type AppConfig struct {
	Repo repository.Repository
}

// Teardown cleans the AppConfig.
//
// This function should be called when the AppConfig is not needed anymore.
func (m *AppConfig) Teardown() {
	m.Repo.Close()
	log.Println("Closed database connection")
}

// InitApp initializes the App configuration variable.
func App() AppConfig {
	once.Do(func() {
		dsn := NewDBOptions(DBName).Dsn()
		app.Repo = db.NewPostgresRepo(dsn)
	})
	return app
}
