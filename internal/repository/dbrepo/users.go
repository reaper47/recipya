package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/contexts"
	"github.com/reaper47/recipya/internal/models"
)

func (m *postgresDBRepo) CreateUser(username, email, password string) (models.User, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var u models.User
	err := m.Pool.QueryRow(ctx, getUserStmt, username, email, -1).Scan(&u.ID, &u.Username, &u.Email, &u.HashedPassword)
	switch err {
	case pgx.ErrNoRows:
		hash, err := auth.HashPassword(password)
		if err != nil {
			return u, fmt.Errorf("could not create user: %s", err)
		}

		var id int64
		err = m.Pool.QueryRow(ctx, insertUserStmt, username, email, hash).Scan(&id)
		if err != nil {
			return u, fmt.Errorf("could not query user row: %s", err)
		}

		u.Populate(id, username, email, string(hash))
		return u, nil
	}
	return u, errors.New("username or email is already taken")
}

func (m *postgresDBRepo) User(id string) models.User {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var u models.User
	_ = m.Pool.QueryRow(ctx, getUserStmt, id, id, -1).Scan(&u.ID, &u.Username, &u.Email, &u.HashedPassword)
	return u
}

func (m *postgresDBRepo) UserByID(id int64) models.User {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var u models.User
	_ = m.Pool.QueryRow(ctx, getUserStmt, "", "", id).Scan(&u.ID, &u.Username, &u.Email, &u.HashedPassword)
	return u
}
