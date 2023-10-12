package services

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services/statements"
	"github.com/reaper47/recipya/internal/units"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"slices"
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	err = goose.SetDialect("sqlite")
	if err != nil {
		panic(err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
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

func (s *SQLiteService) AddCookbook(title string, userID int64) (int64, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var id int64
	err := s.DB.QueryRowContext(ctx, statements.InsertCookbook, title, userID).Scan(&id)
	return id, err
}

func (s *SQLiteService) AddRecipe(r *models.Recipe, userID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var isRecipeExists bool
	err := s.DB.QueryRowContext(ctx, statements.IsRecipeForUserExist, userID, r.Name, r.Description, r.Yield, r.URL).Scan(&isRecipeExists)
	if err != nil {
		return -1, err
	}

	if isRecipeExists {
		return -1, fmt.Errorf("recipe '%s' exists for user %d", r.Name, userID)
	}

	settings, err := s.UserSettings(userID)
	if err != nil {
		return -1, err
	}

	if settings.ConvertAutomatically {
		r, _ = r.ConvertMeasurementSystem(settings.MeasurementSystem)
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// Insert recipe
	var recipeID int64
	err = tx.QueryRowContext(ctx, statements.InsertRecipe, r.Name, r.Description, r.Image, r.Yield, r.URL).Scan(&recipeID)
	if err != nil {
		return -1, err
	}
	r.ID = recipeID

	_, err = tx.ExecContext(ctx, statements.InsertUserRecipe, userID, recipeID)
	if err != nil {
		return -1, err
	}

	// Insert category
	var categoryID int64
	err = tx.QueryRowContext(ctx, statements.InsertCategory, r.Category, userID).Scan(&categoryID)
	if err != nil {
		return -1, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeCategory, categoryID, recipeID)
	if err != nil {
		return -1, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertUserCategory, userID, categoryID)
	if err != nil {
		return -1, err
	}

	// Insert cuisine
	_, err = tx.ExecContext(ctx, statements.InsertCuisine, r.Cuisine, userID)
	if err != nil {
		return -1, err
	}

	var cuisineID int64
	err = tx.QueryRowContext(ctx, statements.SelectCuisineID, r.Cuisine).Scan(&cuisineID)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeCuisine, cuisineID, recipeID)
	if err != nil {
		return -1, err
	}

	// Insert nutrition
	n := r.Nutrition
	_, err = tx.ExecContext(ctx, statements.InsertNutrition, recipeID, n.Calories, n.TotalCarbohydrates, n.Sugars, n.Protein, n.TotalFat, n.SaturatedFat, n.UnsaturatedFat, n.Cholesterol, n.Sodium, n.Fiber)
	if err != nil {
		return -1, err
	}

	// Insert times
	var timesID int64
	err = tx.QueryRowContext(ctx, statements.InsertTimes, int64(r.Times.Prep.Seconds()), int64(r.Times.Cook.Seconds())).Scan(&timesID)
	if err != nil {
		return -1, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeTime, timesID, recipeID)
	if err != nil {
		return -1, err
	}

	// Insert keywords
	for _, keyword := range r.Keywords {
		var keywordID int64
		err := tx.QueryRowContext(ctx, statements.InsertKeyword, keyword).Scan(&keywordID)
		if err != nil {
			return -1, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeKeyword, keywordID, recipeID)
		if err != nil {
			return -1, err
		}
	}

	// Insert instructions
	for i, instruction := range r.Instructions {
		var instructionID int64
		err := tx.QueryRowContext(ctx, statements.InsertInstruction, instruction).Scan(&instructionID)
		if err != nil {
			return -1, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeInstruction, instructionID, recipeID, i)
		if err != nil {
			return -1, err
		}
	}

	// Insert ingredients
	for i, ingredient := range r.Ingredients {
		var ingredientID int64
		ingredient = units.ReplaceDecimalFractions(ingredient)
		err := tx.QueryRowContext(ctx, statements.InsertIngredient, ingredient).Scan(&ingredientID)
		if err != nil {
			return -1, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeIngredient, ingredientID, recipeID, i)
		if err != nil {
			return -1, err
		}
	}

	// Insert tools
	for _, tool := range r.Tools {
		var toolID int64
		err := tx.QueryRowContext(ctx, statements.InsertTool, tool).Scan(&toolID)
		if err != nil {
			return -1, err
		}
		_, err = tx.ExecContext(ctx, statements.InsertRecipeTool, toolID, recipeID)
		if err != nil {
			return -1, err
		}
	}

	var count int64
	err = tx.QueryRowContext(ctx, statements.SelectRecipeCount, userID).Scan(&count)
	return count, tx.Commit()
}

func (s *SQLiteService) AddShareLink(share models.ShareRecipe) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	link := "/r/" + uuid.New().String()
	_, err := s.DB.ExecContext(ctx, statements.InsertShareLink, link, share.RecipeID, share.UserID)
	return link, err
}

func (s *SQLiteService) Categories(userID int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var categories []string
	rows, err := s.DB.QueryContext(ctx, statements.SelectCategories, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c string
		err := rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
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

func (s *SQLiteService) Cookbooks(userID int64) ([]models.Cookbook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbooks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookbooks []models.Cookbook
	for rows.Next() {
		c := models.Cookbook{}
		// TODO: Fetch recipes
		err := rows.Scan(&c.ID, &c.Image, &c.Title, &c.Count)
		if err != nil {
			return nil, err
		}
		cookbooks = append(cookbooks, c)
	}
	return cookbooks, nil
}

func (s *SQLiteService) DeleteAuthToken(userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.DeleteAuthToken, userID)
	return err
}

func (s *SQLiteService) DeleteCookbook(id, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.DeleteCookbook, userID, id)
	return err
}

func (s *SQLiteService) DeleteRecipe(id, userID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	result, err := s.DB.ExecContext(ctx, statements.DeleteRecipe, userID, id)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (s *SQLiteService) GetAuthToken(selector, validator string) (models.AuthToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	token := models.AuthToken{Selector: selector}
	row := s.DB.QueryRowContext(ctx, statements.SelectAuthToken, selector)
	err := row.Scan(&token.ID, &token.HashValidator, &token.Expires, &token.UserID)
	if errors.Is(err, sql.ErrNoRows) {
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

func (s *SQLiteService) IsUserPassword(id int64, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var hash string
	err := s.DB.QueryRowContext(ctx, statements.SelectUserPasswordByID, id).Scan(&hash)
	if err != nil {
		return false
	}

	return auth.VerifyPassword(password, auth.HashedPassword(hash))
}

func (s *SQLiteService) MeasurementSystems(userID int64) ([]units.System, models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		convertAutomatically int64
		groupedSystems       string
		selected             string
	)
	err := s.DB.QueryRowContext(ctx, statements.SelectMeasurementSystems, userID).Scan(&selected, &groupedSystems, &convertAutomatically)
	if err != nil {
		return nil, models.UserSettings{}, err
	}

	systemsStr := strings.Split(groupedSystems, ",")
	systems := make([]units.System, len(systemsStr))
	for i, s := range systemsStr {
		systems[i] = units.NewSystem(s)
	}
	return systems, models.UserSettings{
		ConvertAutomatically: convertAutomatically == 1,
		MeasurementSystem:    units.NewSystem(selected),
	}, nil
}

func (s *SQLiteService) Recipe(id, userID int64) (*models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, statements.SelectRecipe, userID, id)
	r, err := scanRecipe(row)
	if err != nil {
		return nil, err
	}
	return r, tx.Commit()
}

func (s *SQLiteService) Recipes(userID int64) models.Recipes {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectRecipes, userID)
	if err != nil {
		return models.Recipes{}
	}
	defer rows.Close()

	var recipes models.Recipes
	for rows.Next() {
		r, err := scanRecipe(rows)
		if err != nil {
			continue
		}
		recipes = append(recipes, *r)
	}
	return recipes
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRecipe(sc scanner) (*models.Recipe, error) {
	var (
		r            models.Recipe
		ingredients  string
		instructions string
		keywords     string
		tools        string
	)

	err := sc.Scan(
		&r.ID, &r.Name, &r.Description, &r.Image, &r.URL, &r.Yield, &r.CreatedAt, &r.UpdatedAt, &r.Category, &r.Cuisine,
		&ingredients, &instructions, &keywords, &tools, &r.Nutrition.Calories, &r.Nutrition.TotalCarbohydrates,
		&r.Nutrition.Sugars, &r.Nutrition.Protein, &r.Nutrition.TotalFat, &r.Nutrition.SaturatedFat, &r.Nutrition.UnsaturatedFat,
		&r.Nutrition.Cholesterol, &r.Nutrition.Sodium, &r.Nutrition.Fiber, &r.Times.Prep, &r.Times.Cook, &r.Times.Total,
	)
	if err != nil {
		return nil, err
	}

	r.Ingredients = strings.Split(ingredients, "<!---->")
	r.Instructions = strings.Split(instructions, "<!---->")
	r.Keywords = strings.Split(keywords, ",")
	r.Tools = strings.Split(tools, ",")

	r.Times.Prep = r.Times.Prep * time.Second
	r.Times.Cook = r.Times.Cook * time.Second
	r.Times.Total = r.Times.Total * time.Second
	return &r, nil
}

func (s *SQLiteService) RecipeShared(link string) (*models.ShareRecipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var share models.ShareRecipe
	err := s.DB.QueryRowContext(ctx, statements.SelectRecipeShared, link).Scan(&share.RecipeID, &share.UserID)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (s *SQLiteService) RecipeUser(recipeID int64) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var userID int64
	_ = s.DB.QueryRowContext(ctx, statements.SelectRecipeUser, recipeID).Scan(&userID)
	return userID
}

func (s *SQLiteService) Register(email string, hashedPassword auth.HashedPassword) (int64, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var userID int64
	err := s.DB.QueryRowContext(ctx, statements.InsertUser, email, hashedPassword).Scan(&userID)
	return userID, err
}

func (s *SQLiteService) SwitchMeasurementSystem(system units.System, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, statements.UpdateMeasurementSystem, system.String(), userID)
	if err != nil {
		return err
	}

	numConverted := 0
	for _, r := range s.Recipes(userID) {
		converted, err := r.ConvertMeasurementSystem(system)
		if err != nil {
			continue
		}

		_, err = tx.ExecContext(ctx, statements.UpdateRecipeDescription, converted.Description, r.ID)
		if err != nil {
			continue
		}

		for i, ingredient := range converted.Ingredients {
			var ingredientID int64
			err := tx.QueryRowContext(ctx, statements.InsertIngredient, ingredient).Scan(&ingredientID)
			if err != nil {
				continue
			}

			_, err = tx.ExecContext(ctx, statements.UpdateRecipeIngredient, ingredientID, r.Ingredients[i], r.ID)
			if err != nil {
				continue
			}
		}

		for i, instruction := range converted.Instructions {
			var instructionID int64
			err := tx.QueryRowContext(ctx, statements.InsertInstruction, instruction).Scan(&instructionID)
			if err != nil {
				continue
			}

			_, err = tx.ExecContext(ctx, statements.UpdateRecipeInstruction, instructionID, r.Instructions[i], r.ID)
			if err != nil {
				continue
			}
		}

		numConverted++
	}
	return tx.Commit()
}

func (s *SQLiteService) UpdateConvertMeasurementSystem(userID int64, isEnabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateConvertAutomatically, isEnabled, userID)
	return err
}

func (s *SQLiteService) UpdatePassword(userID int64, password auth.HashedPassword) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdatePassword, string(password), userID)
	return err
}

func (s *SQLiteService) UpdateRecipe(updatedRecipe *models.Recipe, userID int64, recipeNum int64) error {
	oldRecipe, err := s.Recipe(recipeNum, userID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recipeID := oldRecipe.ID

	if updatedRecipe.Category != oldRecipe.Category {
		var categoryID int64
		err := tx.QueryRowContext(ctx, statements.InsertCategory, updatedRecipe.Category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, statements.UpdateRecipeCategory, categoryID, recipeID)
		if err != nil {
			return err
		}
	}

	if !slices.Equal(updatedRecipe.Ingredients, oldRecipe.Ingredients) {
		ids := make([]int64, len(updatedRecipe.Ingredients))
		for i, v := range updatedRecipe.Ingredients {
			var id int64
			err := tx.QueryRowContext(ctx, statements.InsertIngredient, v).Scan(&id)
			if err != nil {
				return err
			}
			ids[i] = id
		}

		_, err := tx.ExecContext(ctx, statements.DeleteRecipeIngredients, recipeID)
		if err != nil {
			return err
		}

		for i, id := range ids {
			_, err := tx.ExecContext(ctx, statements.InsertRecipeIngredient, id, recipeID, i)
			if err != nil {
				return err
			}
		}
	}

	if !slices.Equal(updatedRecipe.Instructions, oldRecipe.Instructions) {
		ids := make([]int64, len(updatedRecipe.Instructions))
		for i, v := range updatedRecipe.Instructions {
			var id int64
			err := tx.QueryRowContext(ctx, statements.InsertInstruction, v).Scan(&id)
			if err != nil {
				return err
			}
			ids[i] = id
		}

		_, err := tx.ExecContext(ctx, statements.DeleteRecipeInstructions, recipeID)
		if err != nil {
			return err
		}

		for i, id := range ids {
			_, err := tx.ExecContext(ctx, statements.InsertRecipeInstruction, id, recipeID, i)
			if err != nil {
				return err
			}
		}
	}

	/* TODO: Support editing keywords
	if !slices.Equal(updatedRecipe.Keywords, oldRecipe.Keywords) {
		ids := make([]int64, len(updatedRecipe.Keywords))
		for i, v := range updatedRecipe.Keywords {
			var id int64
			if err := tx.QueryRowContext(ctx, statements.InsertKeyword, v).Scan(&id); err != nil {
				return err
			}
			ids[i] = id
		}

		if _, err := tx.ExecContext(ctx, statements.DeleteRecipeKeywords, recipeID); err != nil {
			return err
		}

		for _, id := range ids {
			if _, err := tx.ExecContext(ctx, statements.InsertRecipeKeyword, id, recipeID); err != nil {
				return err
			}
		}
	}*/

	/* TODO: Support editing tools
	if !slices.Equal(updatedRecipe.Tools, oldRecipe.Tools) {
		ids := make([]int64, len(updatedRecipe.Tools))
		for i, v := range updatedRecipe.Tools {
			var id int64
			if err := tx.QueryRowContext(ctx, statements.InsertTool, v).Scan(&id); err != nil {
				return err
			}
			ids[i] = id
		}

		if _, err := tx.ExecContext(ctx, statements.DeleteRecipeTools, recipeID); err != nil {
			return err
		}

		for _, id := range ids {
			if _, err := tx.ExecContext(ctx, statements.InsertRecipeTool, id, recipeID); err != nil {
				return err
			}
		}
	}*/

	updateFields := make(map[string]any)
	if updatedRecipe.Description != oldRecipe.Description {
		updateFields["description"] = updatedRecipe.Description
	}

	if updatedRecipe.Image != uuid.Nil && updatedRecipe.Image != oldRecipe.Image {
		updateFields["image"] = updatedRecipe.Image.String()
	}

	if updatedRecipe.Name != oldRecipe.Name {
		updateFields["name"] = updatedRecipe.Name
	}

	if updatedRecipe.URL != oldRecipe.URL {
		updateFields["url"] = updatedRecipe.URL
	}

	if updatedRecipe.Yield != oldRecipe.Yield {
		updateFields["yield"] = updatedRecipe.Yield
	}

	fields := []string{"name", "description", "image", "yield", "url"}
	for _, field := range fields {
		if _, ok := updateFields[field]; ok {
			var xs []string
			var args []any
			for _, field := range fields {
				if _, ok := updateFields[field]; ok {
					xs = append(xs, field+" = ?")
					args = append(args, updateFields[field])
				}
			}

			xs[0] = " " + xs[0]
			stmt := "UPDATE recipes SET" + strings.Join(xs, ", ") + " WHERE id = ?"
			args = append(args, recipeID)
			_, err := tx.ExecContext(ctx, stmt, args...)
			if err != nil {
				return err
			}
			break
		}
	}

	if updatedRecipe.Times.Prep != oldRecipe.Times.Prep ||
		updatedRecipe.Times.Cook != oldRecipe.Times.Cook {
		var timesID int64
		err := tx.QueryRowContext(ctx, statements.InsertTimes, int64(updatedRecipe.Times.Prep.Seconds()), int64(updatedRecipe.Times.Cook.Seconds())).Scan(&timesID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, statements.UpdateRecipeTimes, timesID, recipeID)
		if err != nil {
			return err
		}
	}

	if !updatedRecipe.Nutrition.Equal(oldRecipe.Nutrition) {
		var args []any
		var xs []string
		if updatedRecipe.Nutrition.Calories != oldRecipe.Nutrition.Calories {
			xs = append(xs, "calories = ?")
			args = append(args, updatedRecipe.Nutrition.Calories)
		}
		if updatedRecipe.Nutrition.Cholesterol != oldRecipe.Nutrition.Cholesterol {
			xs = append(xs, "cholesterol = ?")
			args = append(args, updatedRecipe.Nutrition.Cholesterol)
		}
		if updatedRecipe.Nutrition.Fiber != oldRecipe.Nutrition.Fiber {
			xs = append(xs, "fiber = ?")
			args = append(args, updatedRecipe.Nutrition.Fiber)
		}
		if updatedRecipe.Nutrition.Protein != oldRecipe.Nutrition.Protein {
			xs = append(xs, "protein = ?")
			args = append(args, updatedRecipe.Nutrition.Protein)
		}
		if updatedRecipe.Nutrition.SaturatedFat != oldRecipe.Nutrition.SaturatedFat {
			xs = append(xs, "saturated_fat = ?")
			args = append(args, updatedRecipe.Nutrition.SaturatedFat)
		}
		if updatedRecipe.Nutrition.Sodium != oldRecipe.Nutrition.Sodium {
			xs = append(xs, "sodium = ?")
			args = append(args, updatedRecipe.Nutrition.Sodium)
		}
		if updatedRecipe.Nutrition.Sugars != oldRecipe.Nutrition.Sugars {
			xs = append(xs, "sugars = ?")
			args = append(args, updatedRecipe.Nutrition.Sugars)
		}
		if updatedRecipe.Nutrition.TotalCarbohydrates != oldRecipe.Nutrition.TotalCarbohydrates {
			xs = append(xs, "total_carbohydrates = ?")
			args = append(args, updatedRecipe.Nutrition.TotalCarbohydrates)
		}
		if updatedRecipe.Nutrition.TotalFat != oldRecipe.Nutrition.TotalFat {
			xs = append(xs, "total_fat = ?")
			args = append(args, updatedRecipe.Nutrition.TotalFat)
		}
		if updatedRecipe.Nutrition.UnsaturatedFat != oldRecipe.Nutrition.UnsaturatedFat {
			xs = append(xs, "unsaturated_fat = ?")
			args = append(args, updatedRecipe.Nutrition.UnsaturatedFat)
		}

		xs[0] = " " + xs[0]
		stmt := "UPDATE nutrition SET" + strings.Join(xs, ", ") + " WHERE recipe_id = ?"
		args = append(args, recipeID)
		_, err := tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteService) UpdateCookbookImage(id int64, image uuid.UUID, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateCookbookImage, image, userID, id)
	return err
}

func (s *SQLiteService) UpdateUserSettingsCookbooksViewMode(userID int64, mode models.ViewMode) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateUserSettingsCookbooksViewMode, mode, userID)
	return err
}

func (s *SQLiteService) UserID(email string) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var id int64
	err := s.DB.QueryRowContext(ctx, statements.SelectUserID, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1
	}
	return id
}

func (s *SQLiteService) UserSettings(userID int64) (models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		convertAutomatically int64
		cookbooksViewMode    int64
		measurementSystem    string
	)
	err := s.DB.QueryRowContext(ctx, statements.SelectUserSettings, userID).Scan(&measurementSystem, &convertAutomatically, &cookbooksViewMode)
	return models.UserSettings{
		CookbooksViewMode:    models.ViewModeFromInt(cookbooksViewMode),
		ConvertAutomatically: convertAutomatically == 1,
		MeasurementSystem:    units.NewSystem(measurementSystem),
	}, err
}

func (s *SQLiteService) UserInitials(userID int64) string {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var email string
	err := s.DB.QueryRowContext(ctx, statements.SelectUserEmail, userID).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
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
		err := rows.Scan(&user.ID, &user.Email)
		if err != nil {
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
	err := s.DB.QueryRowContext(ctx, statements.SelectUserPassword, email).Scan(&id, &hash)
	if err != nil {
		return -1
	}

	if ok := auth.VerifyPassword(password, auth.HashedPassword(hash)); !ok {
		return -1
	}
	return id
}

func (s *SQLiteService) Websites() models.Websites {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var count int64
	err := s.DB.QueryRowContext(ctx, statements.SelectCountWebsites).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
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
		err := rows.Scan(&w.ID, &w.Host, &w.URL)
		if err != nil {
			log.Printf("error scanning website: %q", err)
			continue
		}
		websites[i] = w
		i++
	}
	return websites
}
