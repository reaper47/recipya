package repository_test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/reaper47/recipe-hunter/repository"
)

func TestDb(t *testing.T) {
	t.Run("Init DB with tables", test_InitDbWithTables)
}

func test_InitDbWithTables(t *testing.T) {
	dir, err := ioutil.TempDir("", "sqlite-test-")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(dir, "tmp.db"))
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(dir)
	}()
	repository.InitDb(db)

	rows, err := db.Query("SELECT * FROM sqlite_master WHERE type='table'")
	if err != nil || rows == nil {
		t.Fatal(err)
	}
	defer rows.Close()

	tables := []string{
		"category",
		"ingredient",
		"instruction",
		"nutrition",
		"recipe",
		"recipe_ingredient",
		"recipe_instruction",
		"recipe_tool",
		"tool",
	}

	for rows.Next() {
		var (
			entityType string
			name       string
			tbl        string
			rootpage   int
			query      string
		)
		if err := rows.Scan(&entityType, &name, &tbl, &rootpage, &query); err != nil {
			t.Fatal(err)
		}

		for i, table := range tables {
			if table == name {
				tables = append(tables[:i], tables[i+1:]...)
				break
			}
			t.Fatal("Table " + table + " not in database.")
		}
	}
}

func initTestDb(t *testing.T) (string, *repository.Repository) {
	tmpfile, err := ioutil.TempFile("", "tmp.db")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", tmpfile.Name())
	if err != nil {
		os.Remove(tmpfile.Name())
		t.Fatal(err)
	}

	repository.InitDb(db)
	return tmpfile.Name(), &repository.Repository{DB: db}
}
