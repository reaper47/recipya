package migration

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/migrations"
)

// Up upgrades the database to the next version.
//
// If isAll is true, then all migrations will be applied at once.
// Otherwise, the next migration will be applied.
func Up(db *sql.DB, isAll bool) {
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

	if isAll {
		err = m.Up()
	} else {
		err = m.Steps(1)
	}

	if err != nil {
		log.Fatalln("Error running migration up:", err)
	}

	versionPostUp, _, err := m.Version()
	if err != nil {
		log.Fatalln("Error retrieving database version post-up:", err)
	}

	log.Printf(
		"Updated schema from version %d to version %d successfully.\n",
		versionPreUp,
		versionPostUp,
	)

	err = runPostUpgradeHooks(versionPostUp, db)
	if err != nil {
		log.Fatalln("Error during post upgrade hooks:", err)
	}

	log.Println("Performed data migration successfully.")
}

// Down downgrades the database to the previous version.
func Down(db *sql.DB, isAll bool) {
	m := getInstance(db)

	versionPreDown, _, err := m.Version()
	if err != nil {
		log.Fatalln("Error retrieving database version pre-down:", err)
	}

	if isAll {
		err = m.Down()
	} else {
		err = m.Steps(-1)
	}

	if err != nil {
		log.Fatalln("Error during migration down:", err)
	}

	versionPostDown, _, err := m.Version()
	if err == migrate.ErrNilVersion {
		versionPostDown = 0
	} else if err != nil {
		log.Fatalln("Error retrieving database version post-down:", err)
	}

	log.Printf(
		"Downgrade schema from version %d to version %d successfully.\n",
		versionPreDown,
		versionPostDown,
	)
}

func getInstance(db *sql.DB) *migrate.Migrate {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalln("Unable to create postgres instance:", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalln("Unable to open migrations fs:", err)
	}

	m, err := migrate.NewWithInstance("file", src, config.DBName, driver)
	if err != nil {
		log.Fatalln("Error initializing migration:", err)
	}
	return m
}
