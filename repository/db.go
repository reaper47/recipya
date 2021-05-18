package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/reaper47/recipe-hunter/config"
	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

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

	err = initDb()
	if err != nil {
		log.Fatalf("Failed to initialize the database: '%v'", err)
	}

	log.Printf("Opened database '%v'\n", path)
}

func initDb() error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	for _, t := range allTables {
		_, err = tx.Exec(createTableStmt(t))
		if err != nil {
			log.Fatalf("Could not create table %v: %v", t.name, err)
		}
	}
	return tx.Commit()
}

func createTableStmt(t table) string {
	stmt := "CREATE TABLE IF NOT EXISTS " + t.name + " ("
	var end string
	for key, value := range t.cols {
		if strings.HasPrefix(key, "!") {
			end += value + ", "
		} else {
			stmt += key + " " + value + ", "
		}
	}
	stmt += end + ")"
	return strings.Replace(stmt, ", )", ")", 1)
}
