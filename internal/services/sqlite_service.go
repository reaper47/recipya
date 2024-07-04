package services

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services/statements"
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"github.com/reaper47/recipya/internal/utils/regex"
	_ "modernc.org/sqlite" // Blank import to initialize the SQL driver.
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	shortCtxTimeout  = 3 * time.Second
	longerCtxTimeout = 5 * time.Minute
)

// SQLiteService represents the Service implemented with SQLite.
type SQLiteService struct {
	DB    *sql.DB
	Mutex *sync.Mutex
	FdcDB *sql.DB
}

// NewSQLiteService creates an SQLiteService object.
func NewSQLiteService() *SQLiteService {
	dsnURI := "file:" + filepath.Join(app.DBBasePath, app.RecipyaDB) + "?" +
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
		FdcDB: openFdcDB(),
		Mutex: &sync.Mutex{},
	}
}

func openFdcDB() *sql.DB {
	path := filepath.Join(app.DBBasePath, app.FdcDB)
	db, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

// AddAuthToken adds an authentication token to the database.
func (s *SQLiteService) AddAuthToken(selector, validator string, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.InsertAuthToken, selector, validator, userID)
	return err
}

// AddCookbook adds a cookbook to the database.
func (s *SQLiteService) AddCookbook(title string, userID int64) (int64, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var id int64
	err := s.DB.QueryRowContext(ctx, statements.InsertCookbook, title, uuid.Nil, userID).Scan(&id)
	return id, err
}

// AddCookbookRecipe adds a recipe to the cookbook.
func (s *SQLiteService) AddCookbookRecipe(cookbookID, recipeID, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var exists int64
	err := s.DB.QueryRowContext(ctx, statements.SelectCookbookRecipeExists, cookbookID, userID, recipeID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		return errors.New("recipe or cookbook does not belong to the user")
	}

	_, err = s.DB.ExecContext(ctx, statements.InsertCookbookRecipe, cookbookID, recipeID, cookbookID, userID)
	return err
}

// AddRecipes adds recipes to the user's collection.
// It returns the IDs of these that were successful and the error.
func (s *SQLiteService) AddRecipes(recipes models.Recipes, userID int64, progress chan models.Progress) ([]int64, []models.ReportLog, error) {
	n := len(recipes)
	if n == 0 {
		return nil, nil, errors.New("no recipes to add")
	}

	settings, err := s.UserSettings(userID)
	if err != nil {
		return nil, nil, err
	}

	var (
		errs       []error
		logs       = make([]models.ReportLog, 0, n)
		ids        = make([]int64, 0, n)
		userIDAttr = slog.Int64("userID", userID)
	)

	for i, r := range recipes {
		if progress != nil {
			progress <- models.Progress{Value: i, Total: n}
		}

		id, err := s.addRecipe(r, userID, settings)
		logs = append(logs, models.NewReportLog(r.Name, err))
		if err != nil {
			errs = append(errs, err)
			slog.Error("Skipped recipe", "recipe", r, userIDAttr, "error", err)
			continue
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		if len(errs) > 0 {
			return nil, nil, errs[0]
		}
		return nil, nil, errors.New("no recipes to add")
	}

	s.calculateNutrition(userID, ids, settings, false)
	return ids, logs, nil
}

func (s *SQLiteService) addRecipe(r models.Recipe, userID int64, settings models.UserSettings) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var isRecipeExists bool
	err := s.DB.QueryRowContext(ctx, statements.IsRecipeForUserExist, userID, r.Name, r.Description, r.Yield, r.URL).Scan(&isRecipeExists)
	if err != nil {
		return 0, err
	}

	if isRecipeExists {
		return 0, errors.New("recipe exists")
	}

	if settings.ConvertAutomatically {
		converted, _ := r.ConvertMeasurementSystem(settings.MeasurementSystem)
		if converted != nil {
			r = *converted
		}
	}

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	recipeID, err := s.addRecipeTx(ctx, tx, r, userID)
	if err != nil {
		return 0, err
	}

	return recipeID, tx.Commit()
}

func (s *SQLiteService) addRecipeTx(ctx context.Context, tx *sql.Tx, r models.Recipe, userID int64) (int64, error) {
	// Verification
	if r.Category == "" {
		r.Category = "uncategorized"
	} else {
		delimiters := []string{",", ";", "|"}
		for _, delim := range delimiters {
			if strings.Contains(r.Category, delim) {
				r.Category = strings.Split(r.Category, delim)[0]
			}
		}
	}

	if len(r.Ingredients) == 0 || len(r.Instructions) == 0 {
		return 0, errors.New("missing ingredients or instructions")
	}

	if r.Name == "" {
		return 0, errors.New("missing name of the recipe")
	}

	if r.Yield == 0 {
		r.Yield = 1
	}

	if r.URL == "" {
		r.URL = "Unknown"
	}

	// Insert recipe
	var mainImage uuid.UUID
	if len(r.Images) > 0 {
		mainImage = r.Images[0]
	}

	var recipeID int64
	err := tx.QueryRowContext(ctx, statements.InsertRecipe, r.Name, r.Description, mainImage, r.Yield, r.URL).Scan(&recipeID)
	if err != nil {
		return 0, err
	}
	r.ID = recipeID

	_, err = tx.ExecContext(ctx, statements.InsertUserRecipe, userID, recipeID)
	if err != nil {
		return 0, err
	}

	// Insert additional images
	if len(r.Images) > 1 {
		for _, u := range r.Images[1:] {
			_, err = tx.ExecContext(ctx, statements.InsertAdditionalImageRecipe, recipeID, u)
			if err != nil {
				return 0, err
			}
		}
	}

	// Insert category
	var categoryID int64
	category := r.Category
	before, _, found := strings.Cut(r.Category, ",")
	if found {
		category = before
	}
	err = tx.QueryRowContext(ctx, statements.InsertCategory, category, userID).Scan(&categoryID)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeCategory, categoryID, recipeID)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertUserCategory, userID, categoryID)
	if err != nil {
		return 0, err
	}

	// Insert cuisine
	_, err = tx.ExecContext(ctx, statements.InsertCuisine, r.Cuisine, userID)
	if err != nil {
		return 0, err
	}

	var cuisineID int64
	err = tx.QueryRowContext(ctx, statements.SelectCuisineID, r.Cuisine).Scan(&cuisineID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeCuisine, cuisineID, recipeID)
	if err != nil {
		return 0, err
	}

	// Insert nutrition
	n := r.Nutrition
	n.Clean()
	_, err = tx.ExecContext(ctx, statements.InsertNutrition, recipeID, n.Calories, n.TotalCarbohydrates, n.Sugars, n.Protein, n.TotalFat, n.SaturatedFat, n.UnsaturatedFat, n.TransFat, n.Cholesterol, n.Sodium, n.Fiber, n.IsPerServing)
	if err != nil {
		return 0, err
	}

	// Insert times
	var timesID int64
	err = tx.QueryRowContext(ctx, statements.InsertTimes, int64(r.Times.Prep.Seconds()), int64(r.Times.Cook.Seconds())).Scan(&timesID)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeTime, timesID, recipeID)
	if err != nil {
		return 0, err
	}

	// Insert keywords
	r.Keywords = slices.DeleteFunc(extensions.Unique(r.Keywords), func(s string) bool { return s == "" })
	for _, keyword := range r.Keywords {
		var keywordID int64
		err = tx.QueryRowContext(ctx, statements.InsertKeyword, keyword).Scan(&keywordID)
		if err != nil {
			return 0, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeKeyword, keywordID, recipeID)
		if err != nil {
			return 0, err
		}
	}

	// Insert instructions
	r.Instructions = slices.DeleteFunc(extensions.Unique(r.Instructions), func(s string) bool { return s == "" })
	for i, instruction := range r.Instructions {
		var instructionID int64
		err = tx.QueryRowContext(ctx, statements.InsertInstruction, instruction).Scan(&instructionID)
		if err != nil {
			return 0, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeInstruction, instructionID, recipeID, i)
		if err != nil {
			return 0, err
		}
	}

	// Insert ingredients
	r.Ingredients = slices.DeleteFunc(extensions.Unique(r.Ingredients), func(s string) bool { return s == "" })
	for i, ingredient := range r.Ingredients {
		var ingredientID int64
		ingredient = units.ReplaceDecimalFractions(ingredient)
		err = tx.QueryRowContext(ctx, statements.InsertIngredient, ingredient).Scan(&ingredientID)
		if err != nil {
			return 0, err
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeIngredient, ingredientID, recipeID, i)
		if err != nil {
			return 0, err
		}
	}

	// Insert tools
	r.Tools = slices.DeleteFunc(extensions.Unique(r.Tools), func(t models.HowToItem) bool { return t.Name == "" })
	for i, tool := range r.Tools {
		var toolID int64
		err = tx.QueryRowContext(ctx, statements.InsertTool, tool.Name).Scan(&toolID)
		if err != nil {
			return 0, err
		}

		if tool.Quantity == 0 {
			tool.Quantity = 1
		}

		_, err = tx.ExecContext(ctx, statements.InsertRecipeTool, toolID, recipeID, tool.Quantity, i)
		if err != nil {
			return 0, err
		}
	}

	_, err = tx.ExecContext(ctx, statements.InsertRecipeShadow, recipeID, r.Name, r.Description, r.URL)
	if err != nil {
		return 0, err
	}

	return recipeID, nil
}

// AddRecipeCategory adds a custom recipe category for the user.
func (s *SQLiteService) AddRecipeCategory(name string, userID int64) error {
	// 1. Verify whether category is ok.
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return errors.New("category is invalid")
	}

	categories, err := s.Categories(userID)
	if err != nil {
		return err
	}

	if slices.Contains(categories, name) {
		return errors.New("category is already in use")
	}

	// 2. Add new category
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var categoryID int64
	err = tx.QueryRowContext(ctx, statements.InsertCategory, name).Scan(&categoryID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, statements.InsertUserCategory, userID, categoryID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// AddReport adds a report to the database.
func (s *SQLiteService) AddReport(report models.Report, userID int64) {
	userIDAttr := slog.Int64("userID", userID)
	reportAttr := slog.Any("report", report)

	if len(report.Logs) == 0 {
		slog.Warn("No report to insert into the database", userIDAttr, reportAttr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.Error("AddReport.BeginTx failed", userIDAttr, reportAttr, "error", err)
		return
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, statements.InsertReport, report.Type, report.CreatedAt, report.ExecTime, userID).Scan(&report.ID)
	if err != nil {
		slog.Error("AddReport.InsertReport failed", userIDAttr, reportAttr, "error", err)
		return
	}

	for _, l := range report.Logs {
		_, err = tx.ExecContext(ctx, statements.InsertReportLog, report.ID, l.Title, l.IsSuccess, l.Error)
		if err != nil {
			slog.Error("AddReport.InsertReportLog failed", userIDAttr, reportAttr, "log", l, "error", err)
			continue
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Warn("AddReport.Commit failed", userIDAttr, reportAttr, "error", err)
	}
}

// AddShareLink adds a share link for the recipe.
func (s *SQLiteService) AddShareLink(share models.Share) (string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		stmt string
		link string
		id   int64
	)
	if share.CookbookID > -1 {
		err := s.DB.QueryRowContext(ctx, statements.SelectCookbookSharedLink, share.CookbookID, share.UserID).Scan(&link)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			var exists int64
			err = s.DB.QueryRowContext(ctx, statements.SelectCookbookExists, share.CookbookID, share.UserID).Scan(&exists)
			if err != nil {
				return "", err
			}

			if exists == 0 {
				return "", errors.New("cookbook does not belong to user")
			}

			stmt = statements.InsertShareLinkCookbook
			link = "/c/" + uuid.New().String()
			id = share.CookbookID
		default:
			return link, nil
		}
	} else {
		stmt = statements.InsertShareLink
		link = "/r/" + uuid.New().String()
		id = share.RecipeID
	}

	_, err := s.DB.ExecContext(ctx, stmt, link, id, share.UserID)
	return link, err
}

// AddShareRecipe adds a shared recipe to the user's collection.
func (s *SQLiteService) AddShareRecipe(recipeID, userID int64) (int64, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	var otherUserID int64
	err := s.DB.QueryRowContext(ctx, statements.SelectRecipeSharedFromRecipeID, recipeID).Scan(&otherUserID)
	if err != nil {
		return 0, err
	}

	r, err := s.Recipe(recipeID, otherUserID)
	if err != nil {
		return 0, err
	}

	var isRecipeExists bool
	err = s.DB.QueryRowContext(ctx, statements.IsRecipeForUserExist, userID, r.Name, r.Description, r.Yield, r.URL).Scan(&isRecipeExists)
	if err != nil {
		return 0, err
	}

	if isRecipeExists {
		return 0, errors.New("recipe exists")
	}

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	newRecipeID, err := s.addRecipeTx(ctx, tx, *r, userID)
	if err != nil {
		return 0, err
	}

	return newRecipeID, tx.Commit()
}

// AppInfo gets general information on the application.
func (s *SQLiteService) AppInfo() (models.AppInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		ai              models.AppInfo
		updateAvailable int64
	)

	err := s.DB.QueryRowContext(ctx, statements.SelectAppInfo).Scan(&updateAvailable, &ai.LastUpdatedAt, &ai.LastCheckedUpdateAt)

	ai.IsUpdateAvailable = updateAvailable == 1
	return ai, err
}

// calculateNutrition calculates the nutrition facts for the recipes.
// It is best to run this function in the background because it takes a while per recipe.
func (s *SQLiteService) calculateNutrition(userID int64, recipes []int64, settings models.UserSettings, force bool) {
	if !settings.CalculateNutritionFact {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(recipes))*longerCtxTimeout)
		defer cancel()

		for _, id := range recipes {
			s.Mutex.Lock()
			recipe, err := s.Recipe(id, userID)
			if err != nil {
				slog.Error("CalculateNutrition.Recipe failed", "error", err)
				continue
			}
			s.Mutex.Unlock()

			if !force && !recipe.Nutrition.Equal(models.Nutrition{}) {
				continue
			}

			nutrients, weight, err := s.Nutrients(recipe.Ingredients)
			if err != nil {
				slog.Error("CalculateNutrition.Nutrients failed", "error", err)
				continue
			}

			recipe.Nutrition = nutrients.NutritionFact(weight)
			n := recipe.Nutrition

			s.Mutex.Lock()
			_, err = s.DB.ExecContext(ctx, statements.UpdateNutrition, n.Calories, n.TotalCarbohydrates, n.Sugars, n.Protein, n.TotalFat, n.SaturatedFat, n.UnsaturatedFat, n.Cholesterol, n.Sodium, n.Fiber, n.IsPerServing, id)
			if err != nil {
				slog.Error("CalculateNutrition.UpdateNutrition failed", "error", err)
			}
			s.Mutex.Unlock()
		}
	}()
}

// Categories gets all categories in the database.
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
		err = rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// CheckUpdate checks whether there is a new release for Recipya.
func (s *SQLiteService) CheckUpdate(files FilesService) (models.AppInfo, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	isLatest, _, err := files.IsAppLatest(app.Info.Version)
	if err != nil {
		return models.AppInfo{}, err
	}

	updateAvailable := 0
	if !isLatest {
		updateAvailable = 1
	}

	ai := models.AppInfo{
		IsUpdateAvailable:   !isLatest,
		LastCheckedUpdateAt: time.Now(),
	}
	err = s.DB.QueryRowContext(ctx, statements.UpdateIsUpdateAvailable, updateAvailable).Scan(&ai.LastUpdatedAt, &ai.LastCheckedUpdateAt)
	return ai, err
}

// Confirm confirms the user's account.
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

// Cookbook gets a cookbook belonging to a user.
func (s *SQLiteService) Cookbook(id, userID int64) (models.Cookbook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	var c models.Cookbook
	err := s.DB.QueryRowContext(ctx, statements.SelectCookbook, id, userID).Scan(&c.ID, &c.Title, &c.Image, &c.Count)
	if err != nil {
		return c, err
	}

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbookRecipes, c.ID)
	if err != nil {
		return c, err
	}
	c.Recipes, err = scanRecipes(rows, false)
	return c, err
}

// CookbookByID gets a cookbook by its ID.
func (s *SQLiteService) CookbookByID(id, userID int64) (models.Cookbook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var c models.Cookbook
	err := s.DB.QueryRowContext(ctx, statements.SelectCookbook, id, userID).Scan(&c.ID, &c.Title, &c.Image, &c.Count)
	if err != nil {
		return models.Cookbook{}, err
	}

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbookRecipes, id)
	if err != nil {
		return models.Cookbook{}, err
	}
	defer rows.Close()

	c.Recipes, err = scanRecipes(rows, false)
	return c, err
}

// CookbookRecipe gets a recipe from a cookbook.
func (s *SQLiteService) CookbookRecipe(id, cookbookID int64) (recipe *models.Recipe, userID int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	err = s.DB.QueryRowContext(ctx, statements.SelectCookbookUser, cookbookID).Scan(&userID)
	if err != nil {
		return nil, 0, err
	}

	row := s.DB.QueryRowContext(ctx, statements.SelectCookbookRecipe, cookbookID, id)
	recipe, err = scanRecipe(row, false)
	return recipe, userID, err
}

// CookbookShared checks whether the cookbook is shared.
// It returns a models.Share. Otherwise, an error.
func (s *SQLiteService) CookbookShared(link string) (*models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var share models.Share
	err := s.DB.QueryRowContext(ctx, statements.SelectCookbookShared, link).Scan(&share.CookbookID, &share.UserID)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

// Cookbooks gets a limited number of cookbooks belonging to the user.
func (s *SQLiteService) Cookbooks(userID int64, page uint64) ([]models.Cookbook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbooks, userID, page-1, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookbooks []models.Cookbook
	for rows.Next() {
		var c models.Cookbook
		// TODO: Fetch recipes
		err = rows.Scan(&c.ID, &c.Image, &c.Title, &c.Count)
		if err != nil {
			return nil, err
		}
		cookbooks = append(cookbooks, c)
	}

	return cookbooks, rows.Err()
}

// CookbooksShared gets the user's shared cookbooks.
func (s *SQLiteService) CookbooksShared(userID int64) ([]models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbooksShared, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shares []models.Share
	for rows.Next() {
		share := models.Share{UserID: userID}
		err = rows.Scan(&share.Link, &share.CookbookID)
		if err != nil {
			return shares, err
		}
		shares = append(shares, share)
	}
	return shares, rows.Err()
}

// CookbooksUser gets all the user's cookbooks.
func (s *SQLiteService) CookbooksUser(userID int64) ([]models.Cookbook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectCookbooksUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cookbooks []models.Cookbook
	for rows.Next() {
		var c models.Cookbook
		err = rows.Scan(&c.ID, &c.Title, &c.Image, &c.Count)
		if err != nil {
			return nil, err
		}

		recipeRows, err := s.DB.QueryContext(ctx, statements.SelectCookbookRecipes, c.ID)
		if err != nil {
			return nil, err
		}

		c.Recipes, err = scanRecipes(recipeRows, false)
		if err != nil {
			return nil, err
		}

		cookbooks = append(cookbooks, c)
	}

	return cookbooks, rows.Err()
}

// Counts gets the models.Counts for the user.
func (s *SQLiteService) Counts(userID int64) (models.Counts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var counts models.Counts
	err := s.DB.QueryRowContext(ctx, statements.SelectCounts, userID).Scan(&counts.Cookbooks, &counts.Recipes)
	return counts, err
}

// DeleteAuthToken removes an authentication token from the database.
func (s *SQLiteService) DeleteAuthToken(userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.DeleteAuthToken, userID)
	return err
}

// DeleteRecipeCategory deletes a user's recipe category.
func (s *SQLiteService) DeleteRecipeCategory(name string, userID int64) error {
	if name == "uncategorized" || name == "" {
		return errors.New("category is invalid")
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var categoryID int64
	err = tx.QueryRowContext(ctx, statements.DeleteUserCategory, userID, name).Scan(&categoryID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, statements.UpdateRecipeCategoryReset, categoryID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteCookbook deletes a user's cookbook.
func (s *SQLiteService) DeleteCookbook(id, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.DeleteCookbook, id, userID)
	return err
}

// DeleteRecipe deletes a user's recipe. It returns the number of rows affected.
func (s *SQLiteService) DeleteRecipe(id, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	result, err := s.DB.ExecContext(ctx, statements.DeleteRecipe, userID, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("recipe not found")
	}

	return nil
}

// DeleteRecipeFromCookbook deletes a recipe from a cookbook. It returns the number of recipes in the cookbook.
func (s *SQLiteService) DeleteRecipeFromCookbook(recipeID, cookbookID int64, userID int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.DeleteCookbookRecipe, cookbookID, userID, recipeID)
	if err != nil {
		return -1, err
	}

	var c models.Cookbook
	err = s.DB.QueryRowContext(ctx, statements.SelectCookbook, cookbookID, userID).Scan(&c.ID, &c.Title, &c.Image, &c.Count)
	return c.Count, err
}

// DeleteUser deletes a user and his or her data.
func (s *SQLiteService) DeleteUser(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.DeleteUser, id)
	return err
}

// GetAuthToken gets a non-expired auth token by the selector.
func (s *SQLiteService) GetAuthToken(selector, validator string) (models.AuthToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	token := models.AuthToken{Selector: selector}
	row := s.DB.QueryRowContext(ctx, statements.SelectAuthToken, selector)

	var expiresInt int64
	err := row.Scan(&token.ID, &token.HashValidator, &expiresInt, &token.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.AuthToken{}, err
	}
	token.Expires = time.Unix(expiresInt, 0)

	if auth.DecodeHashValidator(validator) != auth.DecodeHashValidator(token.HashValidator) {
		return models.AuthToken{}, errors.New("unequal hashes")
	}

	return token, nil
}

// Images fetches all distinct image UUIDs for recipes.
// An empty slice is returned when an error occurred.
func (s *SQLiteService) Images() []string {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	xs := make([]string, 0)

	rows, err := s.DB.QueryContext(ctx, statements.SelectDistinctImages)
	if err != nil {
		return xs
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err == nil {
			xs = append(xs, file+app.ImageExt)
		}
	}
	return xs
}

// InitAutologin creates a default user for the autologin feature if no users are present.
func (s *SQLiteService) InitAutologin() error {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	var id int64
	err := s.DB.QueryRowContext(ctx, statements.SelectUserOne).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		hashPassword, err := auth.HashPassword("admin")
		if err != nil {
			return err
		}

		_, err = s.DB.ExecContext(ctx, statements.InsertUser, "admin@autologin.com", hashPassword)
		if err != nil {
			return err
		}

		slog.Info("Created user for autologin with email 'admin@autologin.com' and password 'admin'")
	}

	return nil
}

// IsUserExist checks whether the user is present in the database.
func (s *SQLiteService) IsUserExist(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var exists int64
	_ = s.DB.QueryRowContext(ctx, statements.SelectUserExist, email).Scan(&exists)
	return exists == 1
}

// IsUserPassword checks whether the password is the user's password.
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

// MeasurementSystems gets the units systems, along with the one the user selected, in the database.
func (s *SQLiteService) MeasurementSystems(userID int64) ([]units.System, models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		calculateNutrition   int64
		convertAutomatically int64
		groupedSystems       string
		selected             string
	)
	err := s.DB.QueryRowContext(ctx, statements.SelectMeasurementSystems, userID).Scan(&selected, &groupedSystems, &convertAutomatically, &calculateNutrition)
	if err != nil {
		return nil, models.UserSettings{}, err
	}

	systemsStr := strings.Split(groupedSystems, ",")
	systems := make([]units.System, len(systemsStr))
	for i, s := range systemsStr {
		systems[i] = units.NewSystem(s)
	}
	return systems, models.UserSettings{
		CalculateNutritionFact: calculateNutrition == 1,
		ConvertAutomatically:   convertAutomatically == 1,
		MeasurementSystem:      units.NewSystem(selected),
	}, nil
}

// Nutrients gets the nutrients for the ingredients from the FDC database, along with the total weight.
func (s *SQLiteService) Nutrients(ingredients []string) (models.NutrientsFDC, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(ingredients))
	tokens := make([]units.TokenizedIngredient, len(ingredients))
	for i, ing := range ingredients {
		go func(s string, index int) {
			defer wg.Done()
			tokens[index] = units.NewTokenizedIngredientFromText(s)
		}(ing, i)
	}
	wg.Wait()

	var weight float64
	var nutrients models.NutrientsFDC
	for _, token := range tokens {
		if len(token.Ingredients) == 0 {
			continue
		}

		m, err := token.Measurement.Convert(units.Gram)
		if err != nil {
			m, err = token.Measurement.Convert(units.Millilitre)
			if err == nil {
				weight += m.Quantity
			}
		} else {
			weight += m.Quantity
		}

		stmt := statements.BuildSelectNutrientFDC(token.Ingredients)
		rows, err := s.FdcDB.QueryContext(ctx, stmt)
		if err != nil {
			return nil, 0, err
		}

		for rows.Next() {
			var n models.NutrientFDC
			err = rows.Scan(&n.ID, &n.Name, &n.Amount, &n.UnitName)
			if err != nil {
				return nil, 0, err
			}
			n.Reference = token.Measurement
			nutrients = append(nutrients, n)
		}
		_ = rows.Close()

		err = rows.Err()
		if err != nil {
			return nil, 0, err
		}
	}

	return nutrients, weight, nil
}

// Recipe gets the user's recipe of the given id.
func (s *SQLiteService) Recipe(id, userID int64) (*models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, statements.SelectRecipe, id, userID)
	r, err := scanRecipe(row, false)
	if err != nil {
		return nil, err
	}
	return r, tx.Commit()
}

// Recipes gets the user's recipes.
func (s *SQLiteService) Recipes(userID int64, opts models.SearchOptionsRecipes) models.Recipes {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	params := []any{userID, opts.Page, opts.Page}
	stmt := statements.SelectRecipes
	if !opts.Sort.IsDefault {
		params = []any{userID}

		values := url.Values{}
		values.Add("page", strconv.FormatUint(opts.Page, 10))
		values.Add("sort", opts.Sort.String())

		stmt = statements.BuildSelectPaginatedResults(models.NewSearchOptionsRecipe(values))
		stmt = strings.Replace(stmt, "WHERE recipes.id IN (SELECT id FROM recipes_fts WHERE user_id = ? ORDER BY rank)", "WHERE user_id = ?", 1)
	}

	rows, err := s.DB.QueryContext(ctx, stmt, params...)
	if err != nil {
		return models.Recipes{}
	}

	recipes, err := scanRecipes(rows, true)
	if err != nil {
		return models.Recipes{}
	}

	return recipes
}

// RecipesAll gets all the user's recipes.
func (s *SQLiteService) RecipesAll(userID int64) models.Recipes {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectRecipesAll, userID)
	if err != nil {
		return nil
	}

	recipes, err := scanRecipes(rows, false)
	if err != nil {
		return nil
	}

	return recipes
}

// RecipeShared checks whether the recipe is shared.
// It returns a models.Share. Otherwise, an error.
func (s *SQLiteService) RecipeShared(link string) (*models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var share models.Share
	err := s.DB.QueryRowContext(ctx, statements.SelectRecipeShared, link).Scan(&share.RecipeID, &share.UserID)
	if err != nil {
		return nil, err
	}
	return &share, nil
}

// RecipesShared gets all the user's shared recipes.
func (s *SQLiteService) RecipesShared(userID int64) ([]models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectRecipesShared, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shared []models.Share
	for rows.Next() {
		share := models.Share{UserID: userID}
		err = rows.Scan(&share.Link, &share.RecipeID)
		if err != nil {
			return nil, err
		}
		shared = append(shared, share)
	}
	return shared, rows.Err()
}

// RecipeUser gets the user for which the recipe belongs to.
func (s *SQLiteService) RecipeUser(recipeID int64) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var userID int64
	_ = s.DB.QueryRowContext(ctx, statements.SelectRecipeUser, recipeID).Scan(&userID)
	return userID
}

// Register adds a new user to the store.
func (s *SQLiteService) Register(email string, hashedPassword auth.HashedPassword) (int64, error) {
	if !regex.Email.MatchString(email) || hashedPassword == "" {
		return -1, errors.New("credentials are invalid")
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var userID int64
	err := s.DB.QueryRowContext(ctx, statements.InsertUser, email, hashedPassword).Scan(&userID)
	return userID, err
}

// ReorderCookbookRecipes reorders the recipe indices of a cookbook.
func (s *SQLiteService) ReorderCookbookRecipes(cookbookID int64, recipeIDs []uint64, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, recipeID := range recipeIDs {
		var exists int64
		err = s.DB.QueryRowContext(ctx, statements.SelectCookbookRecipeExists, cookbookID, userID, recipeID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists == 1 {
			_, err = tx.ExecContext(ctx, statements.UpdateCookbookRecipesReorder, i, cookbookID, recipeID)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// Report gets a report of any type belonging to the user.
func (s *SQLiteService) Report(id, userID int64) ([]models.ReportLog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectReport, id, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ReportLog
	for rows.Next() {
		var (
			l         models.ReportLog
			isSuccess int64
		)

		err = rows.Scan(&l.ID, &l.Title, &isSuccess, &l.Error)
		if err != nil {
			return nil, err
		}
		l.IsSuccess = isSuccess == 1

		logs = append(logs, l)
	}

	return logs, nil
}

// ReportsImport gets all import reports.
func (s *SQLiteService) ReportsImport(userID int64) ([]models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, statements.SelectReports, models.ImportReportType, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.Report
	for rows.Next() {
		var (
			report models.Report
			ns     int64
			logIDs string
		)

		err = rows.Scan(&report.ID, &report.CreatedAt, &ns, &logIDs)
		if err != nil {
			return nil, err
		}

		report.ExecTime, err = time.ParseDuration(fmt.Sprintf("%dns", ns))
		if err != nil {
			return nil, err
		}

		report.Logs = make([]models.ReportLog, len(strings.Split(logIDs, ";")))
		reports = append(reports, report)
	}

	slices.SortFunc(reports, func(a, b models.Report) int {
		if b.CreatedAt.Before(a.CreatedAt) {
			return -1
		}
		if a.CreatedAt.After(a.CreatedAt) {
			return +1
		}
		return 0
	})

	return reports, rows.Err()
}

// RestoreUserBackup restores the user's data at the specified date.
func (s *SQLiteService) RestoreUserBackup(backup *models.UserBackup) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, backup.DeleteSQL)
	if err != nil {
		return err
	}

	for _, r := range backup.Recipes {
		_, err = s.addRecipeTx(ctx, tx, r, backup.UserID)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, backup.InsertSQL)
	if err != nil {
		return err
	}

	copyImage := func(name string) error {
		destPath := filepath.Join(app.ImagesDir, name)
		_, err = os.Stat(destPath)
		if err == nil {
			return nil
		}

		src, err := os.Open(filepath.Join(backup.ImagesPath, name))
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.Open(destPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		return err
	}

	files, err := os.ReadDir(backup.ImagesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		err = copyImage(file.Name())
	}

	return tx.Commit()
}

// SearchRecipes searches for recipes based on the configuration.
// It returns the paginated search recipes, the total number of search results and an error.
func (s *SQLiteService) SearchRecipes(opts models.SearchOptionsRecipes, userID int64) (models.Recipes, uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	args := []any{userID}

	arg := opts.Arg()
	if arg != "" {
		var fts string
		if opts.Query != "" {
			fts += opts.Query + "* AND "
		}
		args = append(args, fts+arg)
	} else {
		if opts.Query != "" {
			args = append(args, strings.Join(strings.Fields(opts.Query), " *")+"*")
		}
	}

	if opts.CookbookID > 0 {
		args = append(args, opts.CookbookID)
	}

	rows, err := s.DB.QueryContext(ctx, statements.BuildSelectPaginatedResults(opts), args...)
	if err != nil {
		return nil, 0, err
	}

	var recipes models.Recipes
	for rows.Next() {
		var (
			r        models.Recipe
			img      uuid.UUID
			count    int64
			keywords sql.NullString
		)
		err = rows.Scan(&r.ID, &r.Name, &r.Description, &img, &r.CreatedAt, &r.Category, &keywords, &count)
		if err != nil {
			return models.Recipes{}, 0, err
		}

		if img != uuid.Nil {
			r.Images = []uuid.UUID{img}
		}

		if keywords.Valid && keywords.String != "" {
			xk := strings.Split(keywords.String, ",")
			slices.Sort(xk)
			r.Keywords = xk
		}

		recipes = append(recipes, r)
	}

	var totalCount uint64
	err = s.DB.QueryRowContext(ctx, statements.BuildSelectSearchResultsCount(opts), args...).Scan(&totalCount)
	return recipes, totalCount, err
}

func scanRecipes(rows *sql.Rows, isSearch bool) (models.Recipes, error) {
	defer rows.Close()

	var recipes models.Recipes
	for rows.Next() {
		r, err := scanRecipe(rows, isSearch)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, *r)
	}

	return recipes, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRecipe(sc scanner, isSearch bool) (*models.Recipe, error) {
	var (
		r              = models.NewBaseRecipe()
		mainImage      uuid.UUID
		otherImagesStr string
		ingredients    string
		instructions   string
		isPerServing   int64
		keywords       sql.NullString
		tools          sql.NullString
		count          int64
		err            error
	)

	if isSearch {
		err = sc.Scan(&r.ID, &r.Name, &r.Description, &mainImage, &r.CreatedAt, &r.Category, &keywords, &count)
		if err != nil {
			return nil, err
		}
	} else {
		err = sc.Scan(
			&r.ID, &r.Name, &r.Description, &mainImage, &otherImagesStr, &r.URL, &r.Yield, &r.CreatedAt, &r.UpdatedAt, &r.Category, &r.Cuisine,
			&ingredients, &instructions, &keywords, &tools, &r.Nutrition.Calories, &r.Nutrition.TotalCarbohydrates,
			&r.Nutrition.Sugars, &r.Nutrition.Protein, &r.Nutrition.TotalFat, &r.Nutrition.SaturatedFat, &r.Nutrition.UnsaturatedFat, &r.Nutrition.TransFat,
			&r.Nutrition.Cholesterol, &r.Nutrition.Sodium, &r.Nutrition.Fiber, &isPerServing, &r.Times.Prep, &r.Times.Cook, &r.Times.Total,
			&count,
		)
		if err != nil {
			return nil, err
		}

		r.Ingredients = strings.Split(ingredients, "<!---->")
		r.Instructions = strings.Split(instructions, "<!---->")
		r.Nutrition.IsPerServing = isPerServing == 1

		if tools.Valid {
			parts := strings.Split(tools.String, ",")
			r.Tools = make([]models.HowToItem, 0, len(parts))
			for _, part := range parts {
				var (
					q    int
					name string
				)

				before, after, ok := strings.Cut(part, " ")
				if ok {
					parsed, err := strconv.Atoi(before)
					if err == nil {
						q = parsed
					}

					name = strings.TrimSpace(after)
				}

				r.Tools = append(r.Tools, models.HowToItem{
					Name:     name,
					Quantity: q,
				})
			}
		}

		r.Times.Prep *= time.Second
		r.Times.Cook *= time.Second
		r.Times.Total *= time.Second
	}

	if keywords.Valid && keywords.String != "" {
		xk := strings.Split(keywords.String, ",")
		slices.Sort(xk)
		r.Keywords = xk
	}

	if mainImage != uuid.Nil {
		split := strings.Split(otherImagesStr, ";")
		other := make([]uuid.UUID, 0, len(split))
		for _, s := range split {
			parsed, err := uuid.Parse(s)
			if err != nil {
				if s != "" {
					slog.Warn("Couldn't parse image ID", "recipeID", r.ID, "imageID", s)
				}
				continue
			}
			other = append(other, parsed)
		}
		r.Images = slices.Concat([]uuid.UUID{mainImage}, other)
	}

	return &r, err
}

// SwitchMeasurementSystem sets the user's units system to the desired one.
func (s *SQLiteService) SwitchMeasurementSystem(system units.System, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, statements.UpdateMeasurementSystem, system.String(), userID)
	return err

	// TODO: Figure out what to do with converting all recipes at once.
	/*numConverted := 0
	for _, r := range s.RecipesAll(userID) {
		converted, err := r.ConvertMeasurementSystem(system)
		if err != nil || converted == nil {
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

			_, err = tx.ExecContext(ctx, statements.UpdateRecipeIngredient, ingredientID, converted.Ingredients[i], converted.ID)
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

			_, err = tx.ExecContext(ctx, statements.UpdateRecipeInstruction, instructionID, converted.Instructions[i], converted.ID)
			if err != nil {
				continue
			}
		}

		numConverted++
	}*/
}

// UpdateCalculateNutrition updates the user's calculate nutrition facts automatically setting.
func (s *SQLiteService) UpdateCalculateNutrition(userID int64, isEnabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateCalculateNutrition, isEnabled, userID)
	return err
}

// UpdateConvertMeasurementSystem updates the user's convert automatically setting.
func (s *SQLiteService) UpdateConvertMeasurementSystem(userID int64, isEnabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateConvertAutomatically, isEnabled, userID)
	return err
}

// UpdateCookbookImage updates the image of a user's cookbook.
func (s *SQLiteService) UpdateCookbookImage(id int64, image uuid.UUID, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateCookbookImage, image, userID, id)
	return err
}

// UpdatePassword updates the user's password.
func (s *SQLiteService) UpdatePassword(userID int64, password auth.HashedPassword) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdatePassword, string(password), userID)
	return err
}

// UpdateRecipe updates the recipe with its new values.
func (s *SQLiteService) UpdateRecipe(updatedRecipe *models.Recipe, userID int64, recipeNum int64) error {
	oldRecipe, err := s.Recipe(recipeNum, userID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
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
		if updatedRecipe.Category == "" {
			updatedRecipe.Category = "uncategorized"
		} else {
			delimiters := []string{",", ";", "|"}
			for _, delim := range delimiters {
				if strings.Contains(updatedRecipe.Category, delim) {
					updatedRecipe.Category = strings.Split(updatedRecipe.Category, delim)[0]
				}
			}
		}

		var categoryID int64
		category := updatedRecipe.Category
		before, _, found := strings.Cut(updatedRecipe.Category, ",")
		if found {
			category = before
		}
		err = tx.QueryRowContext(ctx, statements.InsertCategory, category).Scan(&categoryID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, statements.UpdateRecipeCategory, categoryID, recipeID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, statements.InsertUserCategory, userID, categoryID)
		if err != nil {
			return err
		}
	}

	isIngredientsUpdated := !slices.Equal(updatedRecipe.Ingredients, oldRecipe.Ingredients)
	if isIngredientsUpdated {
		updatedRecipe.Ingredients = slices.DeleteFunc(extensions.Unique(updatedRecipe.Ingredients), func(s string) bool { return s == "" })

		if len(updatedRecipe.Ingredients) == 0 {
			return errors.New("missing ingredients")
		}

		ids := make([]int64, 0, len(updatedRecipe.Ingredients))
		for _, v := range updatedRecipe.Ingredients {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertIngredient, v).Scan(&id)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}

		_, err = tx.ExecContext(ctx, statements.DeleteRecipeIngredients, recipeID)
		if err != nil {
			return err
		}

		for i, id := range ids {
			_, err = tx.ExecContext(ctx, statements.InsertRecipeIngredient, id, recipeID, i)
			if err != nil {
				return err
			}
		}
	}

	if !slices.Equal(updatedRecipe.Instructions, oldRecipe.Instructions) {
		updatedRecipe.Instructions = slices.DeleteFunc(extensions.Unique(updatedRecipe.Instructions), func(s string) bool { return s == "" })

		if len(updatedRecipe.Instructions) == 0 {
			return errors.New("missing instructions")
		}

		ids := make([]int64, 0, len(updatedRecipe.Instructions))
		for _, v := range updatedRecipe.Instructions {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertInstruction, v).Scan(&id)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}

		_, err = tx.ExecContext(ctx, statements.DeleteRecipeInstructions, recipeID)
		if err != nil {
			return err
		}

		for i, id := range ids {
			_, err = tx.ExecContext(ctx, statements.InsertRecipeInstruction, id, recipeID, i)
			if err != nil {
				return err
			}
		}
	}

	if !slices.Equal(updatedRecipe.Keywords, oldRecipe.Keywords) {
		updatedRecipe.Keywords = slices.DeleteFunc(updatedRecipe.Keywords, func(s string) bool { return s == "" })

		ids := make([]int64, len(updatedRecipe.Keywords))
		for i, v := range updatedRecipe.Keywords {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertKeyword, v).Scan(&id)
			if err != nil {
				return err
			}
			ids[i] = id
		}

		_, err = tx.ExecContext(ctx, statements.DeleteRecipeKeywords, recipeID)
		if err != nil {
			return err
		}

		for _, id := range ids {
			_, err = tx.ExecContext(ctx, statements.InsertRecipeKeyword, id, recipeID)
			if err != nil {
				return err
			}
		}
	}

	if !slices.Equal(updatedRecipe.Tools, oldRecipe.Tools) {
		updatedRecipe.Tools = slices.DeleteFunc(updatedRecipe.Tools, func(t models.HowToItem) bool { return t.Name == "" })

		ids := make([]int64, 0, len(updatedRecipe.Tools))
		updatedRecipe.Tools = extensions.Unique(updatedRecipe.Tools)
		for _, tool := range updatedRecipe.Tools {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertTool, tool.Name).Scan(&id)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}

		_, err = tx.ExecContext(ctx, statements.DeleteRecipeTools, recipeID)
		if err != nil {
			return err
		}

		for i, id := range ids {
			_, err = tx.ExecContext(ctx, statements.InsertRecipeTool, id, recipeID, updatedRecipe.Tools[i].Quantity, i)
			if err != nil {
				return err
			}
		}
	}

	updateFields := make(map[string]any)
	if updatedRecipe.Description != oldRecipe.Description {
		updateFields["description"] = updatedRecipe.Description
	}

	userIDAttr := slog.Int64("userID", userID)
	recipeIDAttr := slog.Int64("recipeID", recipeID)

	_, err = tx.ExecContext(ctx, statements.DeleteRecipeImages, recipeID, userID)
	if err != nil {
		slog.Error("Failed to delete images.", userIDAttr, recipeIDAttr, "error", err)
		return err
	}

	if len(updatedRecipe.Images) > 0 {
		updateFields["image"] = updatedRecipe.Images[0].String()

		if len(updatedRecipe.Images) > 1 {
			for _, u := range updatedRecipe.Images[1:] {
				_, err = tx.ExecContext(ctx, statements.InsertRecipeImage, recipeID, u)
				if err != nil {
					slog.Warn("Error inserting recipe image", userIDAttr, recipeIDAttr, "image", u, "err", err)
					continue
				}
			}
		}
	}

	if updatedRecipe.Name != oldRecipe.Name {
		if updatedRecipe.Name == "" {
			return errors.New("missing the name of the recipe")
		}
		updateFields["name"] = updatedRecipe.Name
	}

	if updatedRecipe.URL != oldRecipe.URL {
		if updatedRecipe.URL == "" {
			updatedRecipe.URL = "Unknown"
		}
		updateFields["url"] = updatedRecipe.URL
	}

	if updatedRecipe.Yield != oldRecipe.Yield {
		if updatedRecipe.Yield == 0 {
			updatedRecipe.Yield = 1
		}
		updateFields["yield"] = updatedRecipe.Yield
	}

	fields := []string{"name", "description", "image", "yield", "url"}
	for _, field := range fields {
		if _, ok := updateFields[field]; ok {
			var xs []string
			var args []any
			for _, field := range fields {
				if _, ok := updateFields[field]; ok {
					xs = append(xs, field+" = trim(?)")
					args = append(args, updateFields[field])
				}
			}

			if len(xs) == 0 {
				continue
			}

			xs[0] = " " + xs[0]
			stmt := "UPDATE recipes SET" + strings.Join(xs, ", ") + " WHERE id = ?"
			args = append(args, recipeID)
			_, err = tx.ExecContext(ctx, stmt, args...)
			if err != nil {
				return err
			}
			break
		}
	}

	if updatedRecipe.Times.Prep != oldRecipe.Times.Prep ||
		updatedRecipe.Times.Cook != oldRecipe.Times.Cook {
		var timesID int64
		err = tx.QueryRowContext(ctx, statements.InsertTimes, int64(updatedRecipe.Times.Prep.Seconds()), int64(updatedRecipe.Times.Cook.Seconds())).Scan(&timesID)
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
		if updatedRecipe.Nutrition.TotalFat != oldRecipe.Nutrition.TotalFat {
			xs = append(xs, "total_fat = ?")
			args = append(args, updatedRecipe.Nutrition.TotalFat)
		}
		if updatedRecipe.Nutrition.SaturatedFat != oldRecipe.Nutrition.SaturatedFat {
			xs = append(xs, "saturated_fat = ?")
			args = append(args, updatedRecipe.Nutrition.SaturatedFat)
		}
		if updatedRecipe.Nutrition.UnsaturatedFat != oldRecipe.Nutrition.UnsaturatedFat {
			xs = append(xs, "unsaturated_fat = ?")
			args = append(args, updatedRecipe.Nutrition.UnsaturatedFat)
		}
		if updatedRecipe.Nutrition.TransFat != oldRecipe.Nutrition.TransFat {
			xs = append(xs, "trans_fat = ?")
			args = append(args, updatedRecipe.Nutrition.TransFat)
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

		if len(xs) > 0 {
			xs[0] = " " + xs[0]
		}

		stmt := "UPDATE nutrition SET" + strings.Join(xs, ", ") + " WHERE recipe_id = ?"
		args = append(args, recipeID)
		_, err = tx.ExecContext(ctx, stmt, args...)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, statements.UpdateRecipeID, recipeID, recipeID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	if isIngredientsUpdated {
		settings, err := s.UserSettings(userID)
		if err != nil {
			slog.Warn("Could not calculate nutrition", userIDAttr, recipeIDAttr, "error", err)
		} else {
			s.calculateNutrition(userID, []int64{recipeID}, settings, true)
		}
	}

	return nil
}

// UpdateUserSettingsCookbooksViewMode updates the user's preferred cookbooks viewing mode.
func (s *SQLiteService) UpdateUserSettingsCookbooksViewMode(userID int64, mode models.ViewMode) error {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.DB.ExecContext(ctx, statements.UpdateUserSettingsCookbooksViewMode, mode, userID)
	return err
}

// UserID gets the user's id from the email. It returns -1 if user not found.
func (s *SQLiteService) UserID(email string) int64 {
	if !regex.Email.MatchString(email) {
		return -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var id int64
	err := s.DB.QueryRowContext(ctx, statements.SelectUserID, email).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1
	}
	return id
}

// UserSettings gets the user's settings.
func (s *SQLiteService) UserSettings(userID int64) (models.UserSettings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var (
		calculateNutrition   int64
		convertAutomatically int64
		cookbooksViewMode    int64
		measurementSystem    string
	)
	err := s.DB.QueryRowContext(ctx, statements.SelectUserSettings, userID).Scan(&measurementSystem, &convertAutomatically, &cookbooksViewMode, &calculateNutrition)
	return models.UserSettings{
		CalculateNutritionFact: calculateNutrition == 1,
		CookbooksViewMode:      models.ViewModeFromInt(cookbooksViewMode),
		ConvertAutomatically:   convertAutomatically == 1,
		MeasurementSystem:      units.NewSystem(measurementSystem),
	}, err
}

// UserInitials gets the user's initials of maximum two characters.
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

// Users gets all users in the database.
func (s *SQLiteService) Users() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	var users []models.User
	rows, err := s.DB.QueryContext(ctx, statements.SelectUsers)
	if err != nil {
		slog.Error("Failed to fetch users", "error", err)
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Email)
		if err != nil {
			slog.Error("Failed to scan user: %q", "error", err)
			return users
		}
		users = append(users, user)
	}

	_ = rows.Err()
	return users
}

// VerifyLogin checks whether the user provided correct login credentials.
// If yes, their user ID will be returned. Otherwise, -1 is returned.
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

// Websites gets the list of supported websites from which to extract the recipe.
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
		slog.Error("Failed to count websites", "error", err)
		return websites
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var w models.Website
		err = rows.Scan(&w.ID, &w.Host, &w.URL)
		if err != nil {
			slog.Error("Failed to scan website", "error", err)
			continue
		}
		websites[i] = w
		i++
	}

	_ = rows.Err()
	return websites
}
