package services

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/services/statements"
	"github.com/reaper47/recipya/internal/units"
	"io"
	"log"
	_ "modernc.org/sqlite" // Blank import to initialize the SQL driver.
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
	shortCtxTimeout  = 3 * time.Second
	longerCtxTimeout = 1 * time.Minute
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

	// TODO: create default user if autologin enabled

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

// AddRecipe adds a recipe to the user's collection.
func (s *SQLiteService) AddRecipe(r *models.Recipe, userID int64, settings models.UserSettings) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), longerCtxTimeout)
	defer cancel()

	var isRecipeExists bool
	err := s.DB.QueryRowContext(ctx, statements.IsRecipeForUserExist, userID, r.Name, r.Description, r.Yield, r.URL).Scan(&isRecipeExists)
	if err != nil {
		return 0, err
	}

	if isRecipeExists {
		return 0, fmt.Errorf("recipe %q exists for user %d", r.Name, userID)
	}

	if settings.ConvertAutomatically {
		converted, _ := r.ConvertMeasurementSystem(settings.MeasurementSystem)
		if converted != nil {
			r = converted
		}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	recipeID, err := s.AddRecipeTx(ctx, tx, r, userID)
	if err != nil {
		return 0, err
	}

	return recipeID, tx.Commit()
}

// AddRecipeTx adds a recipe to the user's collection using an existing database transaction.
func (s *SQLiteService) AddRecipeTx(ctx context.Context, tx *sql.Tx, r *models.Recipe, userID int64) (int64, error) {
	// Insert recipe
	var recipeID int64
	err := tx.QueryRowContext(ctx, statements.InsertRecipe, r.Name, r.Description, r.Image, r.Yield, r.URL).Scan(&recipeID)
	if err != nil {
		return 0, err
	}
	r.ID = recipeID

	_, err = tx.ExecContext(ctx, statements.InsertUserRecipe, userID, recipeID)
	if err != nil {
		return 0, err
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
	_, err = tx.ExecContext(ctx, statements.InsertNutrition, recipeID, n.Calories, n.TotalCarbohydrates, n.Sugars, n.Protein, n.TotalFat, n.SaturatedFat, n.UnsaturatedFat, n.Cholesterol, n.Sodium, n.Fiber, n.IsPerServing)
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
	for _, tool := range r.Tools {
		var toolID int64
		err = tx.QueryRowContext(ctx, statements.InsertTool, tool).Scan(&toolID)
		if err != nil {
			return 0, err
		}
		_, err = tx.ExecContext(ctx, statements.InsertRecipeTool, toolID, recipeID)
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

// CalculateNutrition calculates the nutrition facts for the recipes.
// It is best to in the background because it takes a while per recipe.
func (s *SQLiteService) CalculateNutrition(userID int64, recipes []int64, settings models.UserSettings) {
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
				log.Printf("CalculateNutrition.Recipe: %q", err)
				continue
			}
			s.Mutex.Unlock()

			if !recipe.Nutrition.Equal(models.Nutrition{}) {
				continue
			}

			nutrients, weight, err := s.Nutrients(recipe.Ingredients)
			if err != nil {
				log.Printf("CalculateNutrition.Nutrients: %q", err)
				continue
			}

			recipe.Nutrition = nutrients.NutritionFact(weight)
			n := recipe.Nutrition

			s.Mutex.Lock()
			_, err = s.DB.ExecContext(ctx, statements.UpdateNutrition, n.Calories, n.TotalCarbohydrates, n.Sugars, n.Protein, n.TotalFat, n.SaturatedFat, n.UnsaturatedFat, n.Cholesterol, n.Sodium, n.Fiber, n.IsPerServing, id)
			if err != nil {
				log.Printf("CalculateNutrition.UpdateNutrition: %q", err)
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
	defer func() {
		_ = rows.Close()
	}()

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
	c.Recipes, err = scanRecipes(rows)
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
	defer func() {
		_ = rows.Close()
	}()

	c.Recipes, err = scanRecipes(rows)
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
	recipe, err = scanRecipe(row)
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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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

		c.Recipes, err = scanRecipes(recipeRows)
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

// DeleteRecipeFromCookbook deletes a recipe from a cookbook. It returns the number of recipes in the cookbook.
func (s *SQLiteService) DeleteRecipeFromCookbook(recipeID, cookbookID uint64, userID int64) (int64, error) {
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
	err := row.Scan(&token.ID, &token.HashValidator, &token.Expires, &token.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.AuthToken{}, err
	}

	if auth.HashValidator(validator) != token.HashValidator {
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
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var file string
		err = rows.Scan(&file)
		if err == nil {
			xs = append(xs, file+".jpg")
		}
	}
	return xs
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
	defer func() {
		_ = tx.Rollback()
	}()

	row := tx.QueryRowContext(ctx, statements.SelectRecipe, id, userID)
	r, err := scanRecipe(row)
	if err != nil {
		return nil, err
	}
	return r, tx.Commit()
}

// Recipes gets the user's recipes.
func (s *SQLiteService) Recipes(userID int64, page uint64) models.Recipes {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	stmt := statements.SelectRecipes

	rows, err := s.DB.QueryContext(ctx, stmt, userID, page, userID)
	if err != nil {
		return models.Recipes{}
	}

	recipes, err := scanRecipes(rows)
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

	recipes, err := scanRecipes(rows)
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
	defer func() {
		_ = rows.Close()
	}()

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
	defer func() {
		_ = tx.Rollback()
	}()

	for i, recipeID := range recipeIDs {
		var exists int64
		err := s.DB.QueryRowContext(ctx, statements.SelectCookbookRecipeExists, cookbookID, userID, recipeID).Scan(&exists)
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
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, backup.DeleteSQL)
	if err != nil {
		return err
	}

	for _, r := range backup.Recipes {
		_, err = s.AddRecipeTx(ctx, tx, &r, backup.UserID)
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
		defer func() {
			_ = src.Close()
		}()

		dest, err := os.Open(destPath)
		if err != nil {
			return err
		}
		defer func() {
			_ = dest.Close()
		}()

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
func (s *SQLiteService) SearchRecipes(query string, options models.SearchOptionsRecipes, userID int64) (models.Recipes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	queries := strings.Split(query, " ")
	stmt := statements.BuildSearchRecipeQuery(queries, options)

	if options.FullSearch {
		for range statements.RecipesFTSFields {
			queries = append(queries, queries...)
		}
	}

	xa := make([]any, len(queries)+1)
	xa[0] = userID
	for i, q := range queries {
		xa[i+1] = q + "*"
	}

	rows, err := s.DB.QueryContext(ctx, stmt, xa...)
	if err != nil {
		return nil, err
	}

	return scanRecipes(rows)
}

func scanRecipes(rows *sql.Rows) (models.Recipes, error) {
	defer func() {
		_ = rows.Close()
	}()

	var recipes models.Recipes
	for rows.Next() {
		r, err := scanRecipe(rows)
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

func scanRecipe(sc scanner) (*models.Recipe, error) {
	var (
		r            models.Recipe
		ingredients  string
		instructions string
		isPerServing int64
		keywords     string
		tools        string
	)

	err := sc.Scan(
		&r.ID, &r.Name, &r.Description, &r.Image, &r.URL, &r.Yield, &r.CreatedAt, &r.UpdatedAt, &r.Category, &r.Cuisine,
		&ingredients, &instructions, &keywords, &tools, &r.Nutrition.Calories, &r.Nutrition.TotalCarbohydrates,
		&r.Nutrition.Sugars, &r.Nutrition.Protein, &r.Nutrition.TotalFat, &r.Nutrition.SaturatedFat, &r.Nutrition.UnsaturatedFat,
		&r.Nutrition.Cholesterol, &r.Nutrition.Sodium, &r.Nutrition.Fiber, &isPerServing, &r.Times.Prep, &r.Times.Cook, &r.Times.Total,
	)
	if err != nil {
		return nil, err
	}

	r.Ingredients = strings.Split(ingredients, "<!---->")
	r.Instructions = strings.Split(instructions, "<!---->")
	r.Keywords = strings.Split(keywords, ",")
	r.Nutrition.IsPerServing = isPerServing == 1
	r.Tools = strings.Split(tools, ",")

	r.Times.Prep *= time.Second
	r.Times.Cook *= time.Second
	r.Times.Total *= time.Second
	return &r, nil
}

// SwitchMeasurementSystem sets the user's units system to the desired one.
func (s *SQLiteService) SwitchMeasurementSystem(system units.System, userID int64) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Hour)
	defer cancel()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, statements.UpdateMeasurementSystem, system.String(), userID)
	if err != nil {
		return err
	}

	numConverted := 0
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
	}
	return tx.Commit()
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

	ctx, cancel := context.WithTimeout(context.Background(), shortCtxTimeout)
	defer cancel()

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	recipeID := oldRecipe.ID

	if updatedRecipe.Category != oldRecipe.Category {
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
	}

	if !slices.Equal(updatedRecipe.Ingredients, oldRecipe.Ingredients) {
		ids := make([]int64, len(updatedRecipe.Ingredients))
		for i, v := range updatedRecipe.Ingredients {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertIngredient, v).Scan(&id)
			if err != nil {
				return err
			}
			ids[i] = id
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
		ids := make([]int64, len(updatedRecipe.Instructions))
		for i, v := range updatedRecipe.Instructions {
			var id int64
			err = tx.QueryRowContext(ctx, statements.InsertInstruction, v).Scan(&id)
			if err != nil {
				return err
			}
			ids[i] = id
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

	return tx.Commit()
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
		log.Printf("error fetching users: %q", err)
		return users
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Email)
		if err != nil {
			log.Printf("error scanning user: %q", err)
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
		log.Printf("websites count error: %q", err)
		return websites
	}
	defer func() {
		_ = rows.Close()
	}()

	i := 0
	for rows.Next() {
		var w models.Website
		err = rows.Scan(&w.ID, &w.Host, &w.URL)
		if err != nil {
			log.Printf("error scanning website: %q", err)
			continue
		}
		websites[i] = w
		i++
	}

	_ = rows.Err()
	return websites
}
