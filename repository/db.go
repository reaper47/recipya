package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/reaper47/recipe-hunter/config"
	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// Db is the database singleton.
func Db() *sql.DB {
	once.Do(createDb)
	return db
}

func createDb() {
	var err error

	path := config.Config.RecipesDb
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %v: %v. Error: %v", dir, dir, err)
	}

	db, err = sql.Open("sqlite", "file:"+path+"?foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = InitDb(db)
	if err != nil {
		log.Fatalf("Failed to initialize the database: '%v'", err)
	}
}

// InitDb initializes the database by creating the tables.
func InitDb(database *sql.DB) error {
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, t := range allTables {
		_, err = tx.Exec(createTableStmt(t))
		if err != nil {
			log.Fatalf("Could not create table %v: %v", t.name, err)
		}
	}
	return tx.Commit()
}
