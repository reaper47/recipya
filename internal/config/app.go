package config

import "github.com/jackc/pgx/v4/pgxpool"

const DBName = "recipya"

type AppConfig struct {
	DB *pgxpool.Pool
}
