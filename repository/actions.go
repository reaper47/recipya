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
func (repo Repository) InsertRecipe(recipe *model.Recipe) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	nutritionID, err := insertNutrition(recipe.Nutrition, tx)
	if err != nil {
		return err
	}

	recipesID, err := insertRecipeBase(recipe, nutritionID, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipesID, recipe.RecipeIngredient, schema.recipesIngredients, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipesID, recipe.RecipeInstructions, schema.recipesInstructions, tx)
	if err != nil {
		return err
	}

	err = insertValues(recipesID, recipe.Tool, schema.recipesTools, tx)
	if err != nil {
		return err
	}

	err = insertCategory(recipesID, recipe.RecipeCategory, tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func insertRecipeBase(recipe *model.Recipe, nutritionID int64, tx *sql.Tx) (int64, error) {
	stmt := "INSERT INTO " + schema.recipes.name + " " +
		"(name, description, recipeYield, nutrition_id, dateCreated, " +
		" url, keywords, prepTime, totalTime, dateModified, image, cookTime) " +
		" VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"

	result, err := tx.Exec(
		stmt,
		strings.ToLower(recipe.Name),
		recipe.Description,
		recipe.RecipeYield,
		nutritionID,
		recipe.DateCreated,
		recipe.Url,
		strings.ToLower(recipe.Keywords),
		recipe.PrepTime,
		recipe.TotalTime,
		recipe.DateModified,
		recipe.Image,
		recipe.CookTime,
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

func insertNutrition(n *model.NutritionSet, tx *sql.Tx) (int64, error) {
	stmt := "SELECT id FROM " + schema.nutrition.name + " " +
		"WHERE protein=? AND sodium=? AND fiber=? AND carbohydrate=? AND fat=? AND " +
		"      saturated_fat=? AND cholesterol=? AND sugar=? AND calories=?"

	var id int64
	row := db.QueryRow(
		stmt,
		n.Protein,
		n.Sodium,
		n.Fiber,
		n.Carbohydrate,
		n.Fat,
		n.SaturatedFat,
		n.Cholesterol,
		n.Sugar,
		n.Calories,
	)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		stmt = "INSERT INTO " + schema.nutrition.name + " " +
			"(protein, sodium, fiber, carbohydrate, fat, " +
			" saturated_fat, cholesterol, sugar, calories) " +
			"VALUES (?,?,?,?,?,?,?,?,?)"

		result, err := tx.Exec(
			stmt,
			n.Protein,
			n.Sodium,
			n.Fiber,
			n.Carbohydrate,
			n.Fat,
			n.SaturatedFat,
			n.Cholesterol,
			n.Sugar,
			n.Calories,
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
		assocTable := strings.Split(t.name, "_")[1]
		stmt := "SELECT id FROM " + assocTable + " WHERE name=\"" + value + "\""
		row := db.QueryRow(stmt)

		var id int64
		err := row.Scan(&id)
		if err == sql.ErrNoRows {
			stmt := "INSERT INTO " + assocTable + " (name) VALUES (?)"
			result, err := tx.Exec(stmt, strings.ToLower(value))
			if err != nil {
				return err
			}

			id, err = result.LastInsertId()
			if err != nil {
				return err
			}
		}

		stmt = "INSERT INTO " + t.name + " (recipes_id, " + assocTable + "_id) VALUES (?,?)"
		_, err = tx.Exec(stmt, recipesID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertCategory(recipesID int64, category string, tx *sql.Tx) error {
	category = strings.ToLower(category)
	stmt := "SELECT id FROM " + schema.categories.name + " WHERE name='" + category + "'"
	row := tx.QueryRow(stmt)

	var categoryID int64
	err := row.Scan(&categoryID)
	if err == sql.ErrNoRows {
		stmt := "INSERT INTO " + schema.categories.name + " (name) VALUES (?)"
		result, err := tx.Exec(stmt, strings.ToLower(category))
		if err != nil {
			return err
		}

		categoryID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	stmt = "UPDATE " + schema.recipes.name + " SET categories_id = ? WHERE id = ?"
	_, err = tx.Exec(stmt, categoryID, recipesID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRecipe updates a recipe in the database if the date modified differs.
func (repo Repository) UpdateRecipe(r *model.Recipe, recipeID int64) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Cateogry
	var categoryID int64
	stmt := "SELECT id FROM " + schema.categories.name + " WHERE name = ?"
	category := strings.ToLower(r.RecipeCategory)
	err = tx.QueryRow(stmt, category).Scan(&categoryID)
	if err == sql.ErrNoRows {
		stmt = "INSERT INTO " + schema.categories.name + " (name) VALUES (?)"
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
	nutritionID, err := insertNutrition(r.Nutrition, tx)
	if err != nil {
		return err
	}

	// Recipe
	stmt = "UPDATE " + schema.recipes.name + " " +
		"SET description = ?, url = ?, image = ?, prepTime = ?, cookTime = ?, totalTime = ?, " +
		"    keywords = ?, recipeYield = ?, categories_id = ?, nutrition_id = ?" +
		"WHERE id = ?"

	_, err = tx.Exec(
		stmt,
		r.Description,
		r.Url,
		r.Image,
		r.PrepTime,
		r.CookTime,
		r.TotalTime,
		r.Keywords,
		r.RecipeYield,
		categoryID,
		nutritionID,
		recipeID,
	)
	if err != nil {
		return err
	}

	// Association tables
	err = updateAssocTable(schema.ingredients, r.RecipeIngredient, recipeID, tx)
	if err != nil {
		return err
	}

	err = updateAssocTable(schema.instructions, r.RecipeInstructions, recipeID, tx)
	if err != nil {
		return err
	}

	err = updateAssocTable(schema.tools, r.Tool, recipeID, tx)
	if err != nil {
		return err
	}

	// No errors, all good so update the dateModified field
	stmt = "UPDATE " + schema.recipes.name + " SET dateModified = ? WHERE id = ?"
	if _, err = tx.Exec(stmt, r.DateModified, recipeID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func updateAssocTable(t table, values []string, recipeID int64, tx *sql.Tx) error {
	stmt := "DELETE FROM " + t.assocTable + " WHERE recipes_id = ?"
	_, err := tx.Exec(stmt, recipeID)
	if err != nil {
		return err
	}

	var ids []int64
	for _, value := range values {
		value = strings.ToLower(value)
		stmt = "SELECT id FROM " + t.name + " WHERE name=\"" + value + "\""
		var id int64
		err = tx.QueryRow(stmt).Scan(&id)
		if err == sql.ErrNoRows {
			stmt = "INSERT INTO " + t.name + " (name) VALUES (?)"
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
		stmt = "INSERT INTO " + t.assocTable + " (recipes_id, " + t.name + "_id) VALUES (?,?)"
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

	r := model.Recipe{Name: name}
	name = strings.ToLower(name)

	fields := "id, description, image, keywords, dateCreated, dateModified, " +
		"	url, prepTime, cookTime, totalTime, recipeYield "

	stmt := "SELECT	" + fields + " " +
		"FROM " + schema.recipes.name + " " +
		"WHERE name = \"" + name + "\""

	err = tx.QueryRow(stmt).Scan(
		&r.ID,
		&r.Description,
		&r.Image,
		&r.Keywords,
		&r.DateCreated,
		&r.DateModified,
		&r.Url,
		&r.PrepTime,
		&r.CookTime,
		&r.TotalTime,
		&r.RecipeYield,
	)
	if err == sql.ErrNoRows {
		tx.Commit()
		return nil, nil
	}

	stmt = "SELECT name FROM " + schema.categories.name + " WHERE id = ?"
	tx.QueryRow(stmt, r.ID).Scan(&r.RecipeCategory)

	ingredients, err := getAssocValues(schema.ingredients, r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.RecipeIngredient = ingredients

	instructions, err := getAssocValues(schema.instructions, r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.RecipeInstructions = instructions

	n, err := getNutritionSet(r.ID, tx)
	if err != nil {
		return nil, err
	}
	r.Nutrition = n

	tools, err := getAssocValues(schema.tools, r.ID, tx)
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
	stmt := "SELECT name " +
		"FROM " + t.name + " " +
		"INNER JOIN " + t.assocTable + " " +
		"ON id = " + t.name + "_id " +
		"WHERE recipes_id = ?"

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
	fields := "carbohydrate, fiber, sodium, sugar, calories, fat, saturated_fat, cholesterol, protein"
	stmt := "SELECT " + fields + " " +
		"FROM " + schema.nutrition.name + " " +
		"WHERE id = (SELECT nutrition_id FROM " + schema.recipes.name + " WHERE id = ?)"

	var n model.NutritionSet
	err := tx.QueryRow(stmt, recipeID).Scan(
		&n.Carbohydrate,
		&n.Fiber,
		&n.Sodium,
		&n.Sugar,
		&n.Calories,
		&n.Fat,
		&n.SaturatedFat,
		&n.Cholesterol,
		&n.Protein,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
