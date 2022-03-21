package config

import (
	"os"
	"testing"
)

func TestDb(t *testing.T) {
	t.Run("DBOptions struct creates the correct DSN", func(t *testing.T) {
		os.Setenv("DB_PROTOCOL", "protocol")
		defer os.Unsetenv("DB_PROTOCOL")
		os.Setenv("DB_HOST", "host")
		defer os.Unsetenv("DB_HOST")
		os.Setenv("DB_PORT", "port")
		defer os.Unsetenv("DB_PORT")
		os.Setenv("DB_USER", "user")
		defer os.Unsetenv("DB_USER")
		os.Setenv("DB_PASSWORD", "password")
		defer os.Unsetenv("DB_PASSWORD")

		opts := NewDBOptions("dbname")
		dsn := opts.Dsn()

		expected := "protocol://user:password@host:port/dbname"
		if dsn != expected {
			t.Fatalf("wanted DSN %s but got %q", expected, dsn)
		}
	})
}
