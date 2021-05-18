package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/reaper47/recipe-hunter/model"
)

const pattern string = "+_+"

type Repository struct {
	DB *sql.DB
}

// InsertRecipe stores a recipe in the database.
func (repo Repository) InsertRecipe(r *model.Recipe) error {
	recipe, err := repo.GetRecipe(r.Name)
	if err != nil {
		return err
	}

	if recipe != nil {
		if recipe.DateCreated == r.DateCreated {
			if recipe.DateModified != r.DateModified {
				return repo.UpdateRecipe(r, recipe.ID)
			}
			return nil
		}
		r.Name = pattern + r.Name
	} else {
		recipe = r
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	recipesID, err := insertRecipeBase(recipe, tx)
	if err != nil {
		return err
	}

	err = insertNutrition(recipesID, recipe.Nutrition, tx)
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

func insertRecipeBase(recipe *model.Recipe, tx *sql.Tx) (int64, error) {
	stmt := "INSERT INTO " + schema.recipes.name + " " +
		"(name, description, recipeYield, dateCreated, " +
		"url, keywords, prepTime, totalTime, dateModified, image, cookTime) " +
		" VALUES (?,?,?,?,?,?,?,?,?,?,?)"

	result, err := tx.Exec(
		stmt,
		strings.ToLower(recipe.Name),
		recipe.Description,
		recipe.RecipeYield,
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

func insertNutrition(id int64, n *model.NutritionSet, tx *sql.Tx) error {
	stmt := "INSERT INTO " + schema.nutrition.name + " " +
		"(recipes_id, protein, sodium, fiber, carbohydrate, fat, " +
		" saturated_fat, cholesterol, sugar, calories) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?)"

	_, err := tx.Exec(
		stmt,
		id,
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
		return err
	}
	return nil
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
func (repo Repository) UpdateRecipe(r *model.Recipe, id int64) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Cateogry
	var categoryID int64 
	stmt := "SELECT id FROM " + schema.categories.name + " WHERE name = ?"
	category := strings.ToLower(r.RecipeCategory)
	 err = tx.QueryRow(stmt, category).Scan(&categoryID);
	 if err != nil {
		return err
	}

	// Recipe
	stmt = "UPDATE " + schema.recipes.name + " " +
		"SET description = ?, url = ?, image = ?, prepTime = ?, cookTime = ?, totalTime = ?, " +
		"    keywords = ?, recipeYield = ?, categories_id = ?" +
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
	)
	if err != nil {
		return err
	}

	// Nutrition

	// Ingredients 

	// Instructions 

	// Tools

	if err = tx.Commit(); err != nil {
		return err
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

	n, err := getNutrition(r.ID, tx)
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

func getNutrition(recipeID int64, tx *sql.Tx) (*model.NutritionSet, error) {
	fields := "carbohydrate, fiber, sodium, sugar, calories, fat, saturated_fat, cholesterol, protein"
	stmt := "SELECT " + fields + " " +
		"FROM " + schema.nutrition.name + " " +
		"WHERE recipes_id = ?"

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
