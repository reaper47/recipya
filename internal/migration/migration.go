package migration

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/reaper47/recipya/internal/config"
)

// Up upgrades the database to the next version.
func Up(db *sql.DB) {
	m := getInstance(db)

	versionPreUp, _, err := m.Version()
	if err == migrate.ErrNilVersion {
		versionPreUp = 0
	} else if err != nil {
		log.Fatalln("Error retrieving database version pre-up:", err)
	}

	err = runPreUpgradeHooks(versionPreUp, db)
	if err != nil {
		log.Fatalln("Error running pre upgrade hooks:", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalln("Error running migration up:", err)
	}

	versionPostUp, _, err := m.Version()
	if err != nil {
		log.Fatalln("Error retrieving database version post-up:", err)
	}

	log.Printf("Updated schema from version %d to version %d successfully.\n", versionPreUp, versionPostUp)

	err = runPostUpgradeHooks(versionPostUp, db)
	if err != nil {
		log.Fatalln("Error during post upgrade hooks:", err)
	}

	log.Println("Performed data migration successfully.")
}

// Down downgrades the database to the previous version.
func Down(db *sql.DB) {
	m := getInstance(db)

	versionPreDown, _, err := m.Version()
	if err != nil {
		log.Fatalln("Error retrieving database version pre-down:", err)
	}

	err = m.Down()
	if err != nil {
		log.Fatalln("Error during migration down:", err)
	}

	versionPostDown, _, err := m.Version()
	if err == migrate.ErrNilVersion {
		versionPostDown = 0
	} else if err != nil {
		log.Fatalln("Error retrieving database version post-down:", err)
	}

	log.Printf("Downgrade schema from version %d to version %d successfully.\n", versionPreDown, versionPostDown)
}

func getInstance(db *sql.DB) *migrate.Migrate {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalln("Unable to create postgres instance:", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get working directory:", err)
	}

	root := strings.SplitN(wd, "internal", 2)[0]
	migrationsDir := filepath.Join(root, "migrations")

	source, err := (&file.File{}).Open("file://" + migrationsDir)
	if err != nil {
		log.Fatalln("Unable to migrations folder:", err)
	}

	m, err := migrate.NewWithInstance("file", source, config.DBName, driver)
	if err != nil {
		log.Fatalln("Error initializing migration:", err)
	}
	return m
}
