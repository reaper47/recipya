package config

import (
	"os"
)

// DBOptions holds connection information for a database.
type DBOptions struct {
	Protocol string
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
}

// Dsn creates the data source name (DSN).
func (d DBOptions) Dsn() string {
	return d.Protocol + "://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.DBName
}

// NewDBOptions creates a new DBOptions struct from environment variables.
//
// The variables are stored in a .env file at the root. The content of this file
// is read by the godotenv package on application startup, which then takes care of
// exporting each entry as an environment variable.
func NewDBOptions(dbname string) DBOptions {
	return DBOptions{
		Protocol: os.Getenv("DB_PROTOCOL"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   dbname,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}
