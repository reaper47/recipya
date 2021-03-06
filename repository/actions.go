package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/reaper47/recipya/consts"
	"github.com/reaper47/recipya/model"
)

// DeleteRecipe deletes a recipe from the database.
func (repo *Repository) DeleteRecipe(id int64) error {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := deleteRecipeStmts()
	var numDeleted int64
	for _, stmt := range stmts {
		result, err := tx.Exec(stmt, id)
		if err != nil {
			return err
		}

		n, err := result.RowsAffected()
		if err != nil {
			return err
		}
		numDeleted += n
	}

	if numDeleted == 0 {
		return consts.ErrEntryNotFound
	}
	return tx.Commit()
}

// GetRecipe gets the recipe from the database that matches the name.
func (repo *Repository) GetRecipe(name string) (*model.Recipe, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	r := model.Recipe{}

	var nutritionID int64
	err = tx.QueryRow(selectRecipeByNameStmt, strings.ToLower(name)).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.Url,
		&r.Image,
		&r.PrepTime,
		&r.CookTime,
		&r.TotalTime,
		&r.RecipeCategory,
		&r.Keywords,
		&r.RecipeYield,
		&nutritionID,
		&r.DateModified,
		&r.DateCreated,
	)
	if err == sql.ErrNoRows {
		tx.Commit()
		return nil, nil
	}

	populateRecipe(&r, tx)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func getAssocValues(t table, recipeID int64, tx *sql.Tx) ([]string, error) {
	stmt := selectAssocValuesStmt(t)
	rows, err := tx.Query(stmt, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		rows.Scan(&value)
		values = append(values, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func getNutritionSet(recipeID int64, tx *sql.Tx) (*model.NutritionSet, error) {
	var n model.NutritionSet
	stmt := selectNutritionSetStmt
	err := tx.QueryRow(stmt, recipeID).Scan(
		&n.Calories,
		&n.Carbohydrate,
		&n.Fat,
		&n.SaturatedFat,
		&n.Cholesterol,
		&n.Protein,
		&n.Sodium,
		&n.Fiber,
		&n.Sugar,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// GetRecipes retrieves all recipes of 0 or 1 category.
func (repo *Repository) GetRecipes(category string, page int, limit int) ([]*model.Recipe, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var stmt string
	if page == -1 {
		stmt = selectRecipesStmt(category)
	} else {
		if category == "" {
			stmt = selectRecipesPageStmt(page, limit)
		} else {
			stmt = selectRecipesByCategoryPageStmt(category, page, limit)
		}
	}

	rows, err := tx.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []*model.Recipe
	for rows.Next() {
		r := model.Recipe{}
		var (
			rowid       int64
			nutritionID int64
		)
		err := rows.Scan(
			&rowid,
			&r.ID,
			&r.Name,
			&r.Description,
			&r.Url,
			&r.Image,
			&r.PrepTime,
			&r.CookTime,
			&r.TotalTime,
			&r.RecipeCategory,
			&r.Keywords,
			&r.RecipeYield,
			&nutritionID,
			&r.DateModified,
			&r.DateCreated,
		)
		if err != nil {
			return nil, err
		}

		if err = populateRecipe(&r, tx); err != nil {
			return nil, err
		}
		recipes = append(recipes, &r)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return recipes, nil
}

// GetRecipesInfo retrieves metadata from the recipes database.
func (repo *Repository) GetRecipesInfo() (*model.RecipesInfo, error) {
	categories, err := repo.GetCategories()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var numRecipes int64
	stmt := selectCountStmt(schema.recipe.name)
	err = tx.QueryRow(stmt).Scan(&numRecipes)
	if err != nil {
		return nil, err
	}

	numRecipesPerCategory := make(map[string]int64)
	for _, category := range categories {
		var count int64
		stmt := selectRecipesCountCategoryStmt(category)
		err = tx.QueryRow(stmt).Scan(&count)
		if err != nil {
			return nil, err
		}
		numRecipesPerCategory[category] = count
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &model.RecipesInfo{
		Total:            numRecipes,
		TotalPerCategory: numRecipesPerCategory,
	}, nil
}

// GetCategories retrieves all recipe categories from the database.
func (repo *Repository) GetCategories() ([]string, error) {
	rows, err := repo.DB.Query(selectCategoriesStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var c string
		err = rows.Scan(&c)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// InsertRecipe stores a recipe in the database.
func (repo *Repository) InsertRecipe(r *model.Recipe) (int64, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	recipeID, err := insertRecipe(r, tx)
	if err != nil {
		return -1, err
	}

	err = insertValues(recipeID, r.RecipeIngredient, schema.recipeIngredient, tx)
	if err != nil {
		return -1, err
	}

	err = insertValues(recipeID, r.RecipeInstructions, schema.recipeInstruction, tx)
	if err != nil {
		return -1, err
	}

	err = insertValues(recipeID, r.Tool, schema.recipeTool, tx)
	if err != nil {
		return -1, err
	}

	return recipeID, tx.Commit()
}

func insertRecipe(r *model.Recipe, tx *sql.Tx) (int64, error) {
	categoryID, err := getCategoryID(r.RecipeCategory, tx)
	if err != nil {
		return -1, err
	}

	nutritionID, err := getNutritionID(r.Nutrition, tx)
	if err != nil {
		return -1, err
	}

	rowIDCategory, err := getRowidCategory(categoryID, tx)
	if err != nil {
		return -1, err
	}

	result, err := tx.Exec(
		insertRecipeStmt,
		rowIDCategory,
		strings.ToLower(r.Name),
		r.Description,
		r.Url,
		r.Image,
		r.PrepTime,
		r.CookTime,
		r.TotalTime,
		categoryID,
		strings.ToLower(r.Keywords),
		r.RecipeYield,
		nutritionID,
		r.DateModified,
		r.DateCreated,
	)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getCategoryID(category string, tx *sql.Tx) (int64, error) {
	category = strings.ToLower(strings.TrimSpace(category))

	var categoryID int64
	err := tx.QueryRow(selectCategoryIdStmt, category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		stmt := insertNameStmt(schema.category.name)
		result, err := tx.Exec(stmt, category)
		if err != nil {
			return -1, err
		}

		categoryID, err = result.LastInsertId()
		if err != nil {
			return -1, err
		}
	}
	return categoryID, nil
}

func getNutritionID(n *model.NutritionSet, tx *sql.Tx) (int64, error) {
	var id int64
	err := db.QueryRow(
		selectNutritionIdStmt,
		n.Calories,
		n.Carbohydrate,
		n.Fat,
		n.SaturatedFat,
		n.Cholesterol,
		n.Protein,
		n.Sodium,
		n.Fiber,
		n.Sugar,
	).Scan(&id)
	if err == sql.ErrNoRows {
		result, err := tx.Exec(
			insertNutritionStmt,
			n.Calories,
			n.Carbohydrate,
			n.Fat,
			n.SaturatedFat,
			n.Cholesterol,
			n.Protein,
			n.Sodium,
			n.Fiber,
			n.Sugar,
		)
		if err != nil {
			return -1, err
		}

		id, err = result.LastInsertId()
		if err != nil {
			return -1, err
		}
	}
	return id, nil
}

func getRowidCategory(categoryID int64, tx *sql.Tx) (int64, error) {
	var id int64
	stmt := selectNextRowidCategoryStmt(categoryID)
	err := tx.QueryRow(stmt).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func insertValues(recipesID int64, values []string, t table, tx *sql.Tx) error {
	for _, value := range values {
		value = strings.ToLower(value)

		var id int64
		stmt := selectIdForNameStmt(t.assocTable)
		err := db.QueryRow(stmt, value).Scan(&id)
		if err == sql.ErrNoRows {
			stmt := insertNameStmt(t.assocTable)
			result, err := tx.Exec(stmt, strings.ToLower(value))
			if err != nil {
				return err
			}

			id, err = result.LastInsertId()
			if err != nil {
				return err
			}
		}

		stmt = insertAssocStmt(t)
		_, err = tx.Exec(stmt, recipesID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateRecipe updates a recipe in the database if the date modified differs.
func (repo *Repository) UpdateRecipe(r *model.Recipe, recipeID int64) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Category
	var categoryID int64
	category := strings.ToLower(r.RecipeCategory)
	err = tx.QueryRow(selectCategoryIdStmt, category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		stmt := insertNameStmt(schema.category.name)
		result, err := tx.Exec(stmt, category)
		if err != nil {
			return err
		}

		categoryID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	// Nutrition
	nutritionID, err := getNutritionID(r.Nutrition, tx)
	if err != nil {
		return err
	}

	// Recipe
	_, err = tx.Exec(
		updateRecipeStmt,
		r.Description,
		r.Url,
		r.Image,
		r.PrepTime,
		r.CookTime,
		r.TotalTime,
		categoryID,
		r.Keywords,
		r.RecipeYield,
		nutritionID,
		recipeID,
	)
	if err != nil {
		return err
	}

	// Association tables
	err = updateAssocTable(schema.ingredient, r.RecipeIngredient, recipeID, tx)
	if err != nil {
		return err
	}

	err = updateAssocTable(schema.instruction, r.RecipeInstructions, recipeID, tx)
	if err != nil {
		return err
	}

	err = updateAssocTable(schema.tool, r.Tool, recipeID, tx)
	if err != nil {
		return err
	}

	// No errors, all good so update the date_modified column
	_, err = tx.Exec(updateDateModifiedStmt, r.DateModified, recipeID)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func updateAssocTable(t table, values []string, recipeID int64, tx *sql.Tx) error {
	stmt := deleteAssocValues(t)
	_, err := tx.Exec(stmt, recipeID)
	if err != nil {
		return err
	}

	var ids []int64
	for _, value := range values {
		value = strings.ToLower(value)

		stmt = selectIdForNameStmt(t.name)
		var id int64
		err = tx.QueryRow(stmt, value).Scan(&id)
		if err == sql.ErrNoRows {
			stmt = insertNameStmt(t.name)
			result, err := tx.Exec(stmt, value)
			if err != nil {
				return err
			}

			id, err = result.LastInsertId()
			if err != nil {
				return err
			}
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		stmt = insertAssocReverseStmt(t)
		_, err = tx.Exec(stmt, recipeID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// ImportRecipe extracts the recipe data from the given URL and
// stores it in the database.
func (repo *Repository) ImportRecipe(url string) (*model.Recipe, error) {
	cmd := exec.Command("python", "./tools/scraper.py", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var recipe model.Recipe
	err = json.Unmarshal(output, &recipe)
	if err != nil {
		return nil, err
	}

	recipeInDb, err := repo.GetRecipe(recipe.Name)
	if err != nil {
		return nil, err
	} else if recipeInDb != nil {
		return recipeInDb, nil
	}

	recipe.ID, err = repo.InsertRecipe(&recipe)
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// SearchMaximizeFridge searches for recipes in the database
// that maximize the number of ingredients taken from the fridge.
func (repo *Repository) SearchMaximizeFridge(ingredients []string, n int) ([]*model.Recipe, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := maximizeFridgeStmt(ingredients, n)
	recipes, err := searchRecipes(stmt, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func searchRecipes(stmt string, tx *sql.Tx) ([]*model.Recipe, error) {
	rows, err := tx.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nutritionID int64
	var recipes []*model.Recipe
	for rows.Next() {
		r := model.Recipe{}
		err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.Url,
			&r.Image,
			&r.PrepTime,
			&r.CookTime,
			&r.TotalTime,
			&r.RecipeCategory,
			&r.Keywords,
			&r.RecipeYield,
			&nutritionID, // &r.Nutrition
			&r.DateModified,
			&r.DateCreated,
		)
		if err != nil {
			return nil, err
		}

		populateRecipe(&r, tx)
		recipes = append(recipes, &r)
	}

	return recipes, nil
}

// SearchMinimizeMissing searches for recipes in the database
// that minimizes the number of ingredients to buy at the store.
func (repo *Repository) SearchMinimizeMissing(
	ingredients []string,
	n int,
) ([]*model.Recipe, error) {
	ctx := context.Background()
	tx, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt := minimizeMissingStmt(ingredients, n)
	recipes, err := searchRecipes(stmt, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func populateRecipe(r *model.Recipe, tx *sql.Tx) error {
	n, err := getNutritionSet(r.ID, tx)
	if err != nil {
		return err
	}
	r.Nutrition = n

	ingredients, err := getAssocValues(schema.ingredient, r.ID, tx)
	if err != nil {
		return err
	}
	r.RecipeIngredient = ingredients

	instructions, err := getAssocValues(schema.instruction, r.ID, tx)
	if err != nil {
		return err
	}
	r.RecipeInstructions = instructions

	tools, err := getAssocValues(schema.tool, r.ID, tx)
	if err != nil {
		return err
	}
	r.Tool = tools

	return nil
}
