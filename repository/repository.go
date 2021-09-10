package repository

import "database/sql"

// Repository holds the database singleton.
type Repository struct {
	DB *sql.DB
}
