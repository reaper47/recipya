package services

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/pressly/goose/v3"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services/statements"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	shortCtxTimeout = 3 * time.Second
)

// SQLiteService represents the Service implemented with SQLite.
type SQLiteService struct {
	DB    *sql.DB
	Mutex *sync.Mutex
}

func NewSQLiteService() *SQLiteService {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(filepath.Dir(exe), "recipya.db")
	dsnURI := "file:" + path + "?" +
		"_pragma=foreign_keys(1)" +
		"&_pragma=journal_mode(wal)" +
		"&_pragma=synchronous(normal)" +
		"&_pragma=temp_store(memory)"

	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	return &SQLiteService{
		DB:    db,
		Mutex: &sync.Mutex{},
	}
}

func (s *SQLiteService) AddAuthToken(selector, validator string, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.InsertAuthToken, selector, validator, userID)
	return err
}

func (s *SQLiteService) Confirm(userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, statements.UpdateIsConfirmed, userID)
	if err != nil {
		return err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *SQLiteService) DeleteAuthToken(userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.DeleteAuthToken, userID)
	return err
}

func (s *SQLiteService) GetAuthToken(selector, validator string) (models.AuthToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	token := models.AuthToken{Selector: selector}
	row := s.DB.QueryRowContext(ctx, statements.SelectAuthToken, selector)
	if err := row.Scan(&token.ID, &token.HashValidator, &token.Expires, &token.UserID); err == sql.ErrNoRows {
		return models.AuthToken{}, err
	}

	if auth.HashValidator(validator) != token.HashValidator {
		return models.AuthToken{}, errors.New("unequal hashes")
	}

	return token, nil
}

func (s *SQLiteService) IsUserExist(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var exists int64
	_ = s.DB.QueryRowContext(ctx, statements.SelectUserExist, email).Scan(&exists)
	return exists == 1
}

func (s *SQLiteService) Register(email string, hashedPassword auth.HashedPassword) (int64, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, statements.InsertUser, email, hashedPassword)
	if err != nil {
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return userID, err
}

func (s *SQLiteService) UserInitials(userID int64) string {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var email string
	if err := s.DB.QueryRowContext(ctx, statements.SelectUserEmail, userID).Scan(&email); err == sql.ErrNoRows {
		return ""
	}
	return string(strings.ToUpper(email)[0])
}

func (s *SQLiteService) Users() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var users []models.User
	rows, err := s.DB.QueryContext(ctx, statements.SelectUsers)
	if err != nil {
		log.Printf("error fetching users: %q", err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			log.Printf("error scanning user: %q", err)
			return users
		}
		users = append(users, user)
	}
	return users
}

func (s *SQLiteService) VerifyLogin(email, password string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		id   int64
		hash string
	)
	if err := s.DB.QueryRowContext(ctx, statements.SelectUserPassword, email).Scan(&id, &hash); err != nil {
		return -1
	}

	if isOk := auth.VerifyPassword(password, auth.HashedPassword(hash)); !isOk {
		return -1
	}
	return id
}

func (s *SQLiteService) Websites() models.Websites {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var count int64
	if err := s.DB.QueryRowContext(ctx, statements.SelectCountWebsites).Scan(&count); err == sql.ErrNoRows {
		return models.Websites{}
	}

	websites := make(models.Websites, count)
	rows, err := s.DB.QueryContext(ctx, statements.SelectWebsites)
	if err != nil {
		log.Printf("websites count error: %q", err)
		return websites
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var w models.Website
		if err := rows.Scan(&w.ID, &w.Host, &w.URL); err != nil {
			log.Printf("error scanning website: %q", err)
			continue
		}
		websites[i] = w
		i++
	}
	return websites
}

func (s *SQLiteService) WebsitesSearch(query string) models.Websites {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	websites := make(models.Websites, 0)

	rows, err := s.DB.QueryContext(ctx, statements.SelectWebsitesSearch, "%"+query+"%")
	if err != nil {
		log.Printf("websites search error: %q", err)
		return websites
	}
	defer rows.Close()

	for rows.Next() {
		var w models.Website
		if err := rows.Scan(&w.ID, &w.Host, &w.URL); err != nil {
			log.Printf("error scanning searched website: %q", err)
			continue
		}
		websites = append(websites, w)
	}

	return websites
}
