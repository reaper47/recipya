package config

import (
	"os"
	"testing"
)

func TestConfigDb(t *testing.T) {
	t.Run("DBOptions struct creates the correct DSN", func(t *testing.T) {
		_ = os.Setenv("DB_PROTOCOL", "protocol")
		_ = os.Setenv("DB_HOST", "host")
		_ = os.Setenv("DB_PORT", "port")
		_ = os.Setenv("DB_USER", "user")
		_ = os.Setenv("DB_PASSWORD", "password")
		defer func() {
			_ = os.Unsetenv("DB_PROTOCOL")
			_ = os.Unsetenv("DB_HOST")
			_ = os.Unsetenv("DB_PORT")
			_ = os.Unsetenv("DB_USER")
			_ = os.Unsetenv("DB_PASSWORD")
		}()

		opts := NewDBOptions("dbname")
		dsn := opts.Dsn()

		expected := "protocol://user:password@host:port/dbname"
		if dsn != expected {
			t.Fatalf("wanted DSN %s but got %q", expected, dsn)
		}
	})
}
