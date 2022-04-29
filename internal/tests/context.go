package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/contexts"
	"github.com/reaper47/recipya/internal/driver"
	"github.com/reaper47/recipya/internal/migration"
)

const testDbName = "recipya-test"

// Context holds a contextual information on a connection.
type Context struct {
	roleName string
	pool     *pgxpool.Pool
}

// Reset resets the database.
func (c Context) Reset() {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	_, err := c.pool.Exec(ctx, "DELETE FROM recipes")
	if err != nil {
		log.Fatalln("Could not reset database:", err)
	}
}

// Close closes the connection pool.
func (c Context) Close() {
	c.pool.Close()

	// Reconnect as our root user
	dbOptions := config.NewDBOptions(testDbName)
	pool := driver.ConnectPostgres(dbOptions.Dsn())

	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	dropSchema(ctx, c.roleName, pool)
	dropRole(ctx, c.roleName, pool, "Close")

	pool.Close()
}

// NewContext creates a new context.
func NewContext() Context {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	dbOptions := config.NewDBOptions(testDbName)
	pool := driver.ConnectPostgres(dbOptions.Dsn())
	roleName := createRole(ctx, pool)

	query := fmt.Sprintf("CREATE SCHEMA %s AUTHORIZATION %s", roleName, roleName)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		dropRole(ctx, roleName, pool, "NewContext")
		log.Fatalf("Could not create schema '%s' for tests: %s\n", roleName, err)
	}

	pool.Close()

	dbOptions.User = roleName
	dbOptions.Password = roleName
	dsn := dbOptions.Dsn()
	db := driver.ConnectSqlDB(dsn)
	migration.Up(db, true)
	err = db.Close()
	if err != nil {
		dropRole(ctx, roleName, pool, "NewContext")
		dropSchema(ctx, roleName, pool)
		log.Fatalln("Could not close database:", err)
	}

	pool = driver.ConnectPostgres(dsn)

	return Context{
		roleName: roleName,
		pool:     pool,
	}
}

func createRole(ctx context.Context, pool *pgxpool.Pool) string {
	token := make([]byte, 4)
	rand.Read(token)
	roleName := fmt.Sprintf("a%x", token)

	query := fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s'", roleName, roleName)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Could not create role '%s' for tests: %s\n", roleName, err)
	}
	return roleName
}

func dropSchema(ctx context.Context, roleName string, pool *pgxpool.Pool) {
	query := fmt.Sprintf("DROP SCHEMA %s CASCADE", roleName)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Could not drop schema '%s' when closing database: %s\n", roleName, err)
	}
}

func dropRole(ctx context.Context, roleName string, pool *pgxpool.Pool, from string) {
	query := fmt.Sprintf("DROP ROLE %s", roleName)
	_, err := pool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("%s: could not drop role '%s' when closing database: %s\n", from, roleName, err)
	}
}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}

	root := strings.SplitN(dir, "internal", 2)[0]
	env := fmt.Sprintf("%s/.env", root)

	err = godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}
