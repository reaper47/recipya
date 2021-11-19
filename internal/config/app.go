package config

import "github.com/jackc/pgx/v4/pgxpool"

// DBName is name the database.
const DBName = "recipya"

// AppConfig holds configuration data for the application.
type AppConfig struct {
	DB *pgxpool.Pool
}
