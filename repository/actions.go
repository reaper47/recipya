package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/reaper47/recipe-hunter/model"
)

type Repository struct {
	DB *sql.DB
}

// InsertRecipe stores a recipe in the database.
func (repo Repository) InsertRecipe(r *model.Recipe) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	categoryID, err := getCategoryID(r.RecipeCategory, tx)
	if err != nil {
		return err
	}

	nutritionID, err := getNutritionID(r.Nutrition, tx)
	if err != nil {
		return err
	}

	recipeID, err := insertRecipe(r, categoryID, nutritionID, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipeID, r.RecipeIngredient, schema.recipeIngredient, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipeID, r.RecipeInstructions, schema.recipeInstruction, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipeID, r.Tool, schema.recipeTool, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func insertRecipe(
	r *model.Recipe,
	categoryID int64,
	nutritionID int64,
	tx *sql.Tx,
) (int64, error) {
	result, err := tx.Exec(
		insertRecipeStmt,
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

func insertValues(
	recipesID int64,
	values []string,
	t table,
	tx *sql.Tx,
) error {
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

func getCategoryID(category string, tx *sql.Tx) (int64, error) {
	category = strings.ToLower(category)

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

// UpdateRecipe updates a recipe in the database if the date modified differs.
func (repo Repository) UpdateRecipe(r *model.Recipe, recipeID int64) error {
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

// GetRecipe gets the recipe from the database that matches the name.
func (repo Repository) GetRecipe(name string) (*model.Recipe, error) {
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
		&nutritionID, // &r.Nutrition
		&r.DateModified,
		&r.DateCreated,
	)
	if err == sql.ErrNoRows {
		tx.Commit()
		return nil, nil
	}

	n, err := getNutritionSet(r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.Nutrition = n

	ingredients, err := getAssocValues(schema.ingredient, r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.RecipeIngredient = ingredients

	instructions, err := getAssocValues(schema.instruction, r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.RecipeInstructions = instructions

	tools, err := getAssocValues(schema.tool, r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.Tool = tools

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
