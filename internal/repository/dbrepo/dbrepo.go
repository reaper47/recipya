package db

import "github.com/jackc/pgx/v4/pgxpool"

type postgresDBRepo struct {
	Pool *pgxpool.Pool
}
