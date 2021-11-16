package migration

import (
	"database/sql"
)

func runPreUpgradeHooks(version uint, db *sql.DB) error {
	return nil
}

func runPostUpgradeHooks(version uint, db *sql.DB) error {
	return nil
}
