package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/reaper47/recipya/internal/driver"
	"github.com/reaper47/recipya/internal/repository"
)

type postgresDBRepo struct {
	Pool *pgxpool.Pool
}

// NewPostgresRepo creates a new PostgreSQL repository.
func NewPostgresRepo(dsn string) repository.Repository {
	return &postgresDBRepo{
		Pool: driver.ConnectPostgres(dsn),
	}
}

func (m *postgresDBRepo) Close() {
	m.Pool.Close()
}
