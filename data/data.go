package data

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

// PopulateBlacklistIngredients adds the blacklisted ingredients
// listed under the 'data/blacklist_ingredients.txt' file to the database.
func PopulateBlacklistIngredients(stmt string, db *sql.DB) error {
	return populate("data/blacklist_units.txt", stmt, db)
}

// PopulateFruitsVeggies adds the produce under the
// 'data/fruits_veggies.txt' file to the database.
func PopulateFruitsVeggies(stmt string, db *sql.DB) error {
	return populate("data/fruits_veggies.txt", stmt, db)
}

func populate(fname string, stmt string, db *sql.DB) error {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Read %s error: %v\n", fname, err)
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var items []interface{}
	var values []string
	for scanner.Scan() {
		values = append(values, "(?)")
		items = append(items, scanner.Text())
	}

	_, err = db.Exec(fmt.Sprintf("%s %s", stmt, strings.Join(values, ",")), items...)
	return err
}
