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

// CreateUsers stores a new user in the database.
func (m *postgresDBRepo) CreateUser(username, email, password string) (models.User, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var u models.User
	err := m.Pool.QueryRow(ctx, getUserStmt, username, email).Scan(&u.ID, &u.Username, &u.Email, &u.HashedPassword)
	switch err {
	case pgx.ErrNoRows:
		hash, err := auth.HashPassword(password)
		if err != nil {
			return u, fmt.Errorf("could not create user: %s", err)
		}

		_, err = m.Pool.Exec(ctx, insertUserStmt, username, email, hash)
		if err != nil {
			return u, err
		}
		return u, nil
	}
	return u, errors.New("username or email is already taken")
}

// GetUser gets a user from the database based on the username or email.
func (m *postgresDBRepo) GetUser(id string) models.User {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var u models.User
	_ = m.Pool.QueryRow(ctx, getUserStmt, id, id).Scan(&u.ID, &u.Username, &u.Email, &u.HashedPassword)
	return u
}
