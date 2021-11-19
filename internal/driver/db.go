package driver

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reaper47/recipya/internal/contexts"
)

// ConnectPostgres creates a pgx connection pool from the PostgreSQL DSN.
func ConnectPostgres(dsn string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}

	ctx, cancel := contexts.DBContext()
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping the database (%s): %s\n", dsn, err)
	}
	return pool
}

// ConnectSqlDB creates a connection from the PostgreSQL DSN.
func ConnectSqlDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}

	ctx, cancel := contexts.DBContext()
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalf("Unable to ping the database (%s): %s\n", dsn, err)
	}
	return db
}
