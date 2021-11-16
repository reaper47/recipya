package db

import "github.com/jackc/pgx/v4/pgxpool"

type posgresDBRepo struct {
	Pool *pgxpool.Pool
}
