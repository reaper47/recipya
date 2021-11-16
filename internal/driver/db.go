package driver

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectPostgres() *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		log.Fatalln("unable to connect to database: ", err)
	}
	return pool
}
