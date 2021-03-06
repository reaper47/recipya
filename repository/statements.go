package repository

import (
	"strconv"
	"strings"
)

func createTableStmt(t table) string {
	stmt := "CREATE TABLE IF NOT EXISTS " + t.name + " ("
	var end string
	for key, value := range t.cols {
		if strings.HasPrefix(key, "!") {
			end += value + ", "
		} else {
			stmt += key + " " + value + ", "
		}
	}
	stmt += end + ")"
	return strings.Replace(stmt, ", )", ")", 1)
}

// SELECT statements
const recipeFields = "id, name, description, url, image, prep_time, cook_time, total_time, " +
	"(SELECT name FROM category WHERE id=category_id) AS category, " +
	"keywords, yield, nutrition_id, date_modified, date_created"

var selectRecipeByNameStmt = "SELECT " + recipeFields + " " +
	"FROM " + schema.recipe.name + " " +
	"WHERE name=?"

func selectRecipesStmt(category string) string {
	recipes := "SELECT rowid_category," + recipeFields + " FROM " + schema.recipe.name
	if category != "" {
		recipes += " WHERE category_id=(SELECT id FROM category WHERE name='" + category + "')"
	}
	return recipes
}

func selectRecipesPageStmt(page int, limit int) string {
	recipes := "SELECT rowid," + recipeFields + " FROM " + schema.recipe.name

	var condition string
	if page == 1 {
		condition = "rowid >= 0"
	} else {
		condition = "rowid > " + strconv.Itoa(page*limit-limit)
	}

	return "SELECT * FROM (" + recipes + ") " +
		"WHERE " + condition + " " +
		"ORDER BY rowid " +
		"LIMIT " + strconv.Itoa(limit)
}

func selectRecipesByCategoryPageStmt(category string, page int, limit int) string {
	return "SELECT * FROM (" + selectRecipesStmt(category) + ") " +
		"WHERE rowid_category >= " + strconv.Itoa(page*limit-limit) + " " +
		"ORDER BY rowid_category " +
		"LIMIT " + strconv.Itoa(limit)
}

var selectCategoryIdStmt = "SELECT id " +
	"FROM " + schema.category.name + " " +
	"WHERE name=?"

var selectCategoriesStmt = "SELECT name " +
	"FROM " + schema.category.name + " " +
	"ORDER BY name ASC"

var selectNutritionIdStmt = "SELECT id FROM " + schema.nutrition.name + " " +
	"WHERE calories=? AND carbohydrate=? AND fat=? AND " +
	"	   saturated_fat=? AND cholesterol=? AND protein=? AND " +
	"	   sodium=? AND fiber=? AND sugar=?"

var selectNutritionSetStmt = "SELECT calories, carbohydrate, fat, saturated_fat, " +
	"       cholesterol, protein, sodium, fiber, sugar " +
	"FROM " + schema.nutrition.name + " " +
	"WHERE id=(" +
	"	SELECT nutrition_id " +
	"	FROM " + schema.recipe.name +
	"	WHERE id=?" +
	")"

func selectCountStmt(tname string) string {
	return "SELECT COUNT(*) FROM " + tname
}

func selectColStmt(col, tname string) string {
	return "SELECT " + col + " FROM " + tname
}

func selectRecipesCountCategoryStmt(category string) string {
	return "SELECT COUNT(*) " +
		"FROM " + schema.recipe.name + " " +
		"WHERE category_id = (" +
		"	SELECT id " +
		"	FROM " + schema.category.name + " " +
		"	WHERE name = '" + category + "'" +
		")"
}

func maximizeFridgeStmt(ingredients []string, limit int) string {
	return createSearchStmt(ingredients, limit, false)
}

func minimizeMissingStmt(ingredients []string, limit int) string {
	return createSearchStmt(ingredients, limit, true)
}

func createSearchStmt(ingredients []string, limit int, isBuyMinIngredients bool) string {
	selectNumIngredients := "(" +
		"SELECT COUNT(*) " +
		"	FROM (" +
		"		SELECT name " +
		"		FROM " + schema.ingredient.name + " " +
		"		INNER JOIN " + schema.ingredient.assocTable + " " +
		"		ON id=" + schema.ingredient.assocTable + ".ingredient_id " +
		"		WHERE " + schema.ingredient.assocTable + ".recipe_id=" + schema.recipe.name + ".id" +
		"	) " +
		"	WHERE (" + createLikeStmt(ingredients) + ")" +
		") "

	selectTotalIngredients := "SELECT COUNT(*) " +
		"FROM " + schema.ingredient.name + " " +
		"INNER JOIN " + schema.ingredient.assocTable + " " +
		"ON id=" + schema.ingredient.assocTable + ".ingredient_id " +
		"WHERE " + schema.ingredient.assocTable + ".recipe_id=" + schema.recipe.name + ".id"

	return "SELECT " + recipeFields +
		"	FROM (" +
		"		SELECT " + schema.recipe.name + ".id, " +
		"		(" + selectNumIngredients + ") AS num_ingredients, " +
		"		(" + selectTotalIngredients + ") AS total_ingredients, " +
		"		recipe.* " +
		"	FROM " + schema.recipe.name + " " +
		"   WHERE num_ingredients >= 1 " +
		" " + orderByMode(isBuyMinIngredients) + " " +
		"	LIMIT " + strconv.Itoa(limit) +
		")"
}

func orderByMode(isMinimizeMissing bool) string {
	if isMinimizeMissing {
		return "ORDER BY (total_ingredients - num_ingredients) ASC, total_ingredients"
	}
	return "ORDER BY num_ingredients DESC, total_ingredients"
}

func createLikeStmt(values []string) string {
	stmt := ""
	for i, value := range values {
		stmt += "name LIKE '%" + value + "%' "
		if i != len(values)-1 {
			stmt += "OR "
		}
	}
	return stmt
}

func selectAssocValuesStmt(t table) string {
	return "SELECT name " +
		"FROM " + t.name + " " +
		"INNER JOIN " + t.assocTable + " " +
		"ON id=" + t.name + "_id " +
		"WHERE recipe_id=?"
}

func selectIdForNameStmt(tname string) string {
	return "SELECT id " +
		"FROM " + tname + " " +
		"WHERE name=?"
}

func selectNextRowidCategoryStmt(id int64) string {
	stmt := "SELECT rowid_category " +
		"FROM " + schema.recipe.name + " " +
		"WHERE category_id=" + strconv.FormatInt(id, 10) + " " +
		"ORDER BY rowid_category DESC " +
		"LIMIT 1"
	return "SELECT IIF((" + stmt + ") IS NOT NULL, (" + stmt + ") + 1, 0)"
}

// INSERT statements
var insertRecipeStmt = "INSERT INTO " + schema.recipe.name + " (" +
	"rowid_category, name, description, url, image, prep_time, " +
	"cook_time, total_time, category_id, keywords, yield, " +
	"nutrition_id, date_modified, date_created" +
	") VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

var insertNutritionStmt = "INSERT INTO " + schema.nutrition.name + " (" +
	"calories, carbohydrate, fat, saturated_fat, " +
	"cholesterol, protein, sodium, fiber, sugar" +
	") VALUES (?,?,?,?,?,?,?,?,?)"

func insertNamesStmt(tname string) string {
	return "INSERT INTO " + tname + " (name) VALUES"
}

func insertNameStmt(tname string) string {
	return "INSERT INTO " + tname + " (" +
		"name" +
		") VALUES (?)"
}

func insertAssocStmt(t table) string {
	return "INSERT INTO " + t.name + " (" +
		"recipe_id, " + t.assocTable + "_id" +
		") VALUES (?,?)"

}

func insertAssocReverseStmt(t table) string {
	return "INSERT INTO " + t.assocTable + " (" +
		"recipe_id, " + t.name + "_id" +
		") VALUES (?,?)"
}

// UPDATE statements
var updateRecipeStmt = "UPDATE " + schema.recipe.name + " " +
	"SET description=?, url=?, image=?, prep_time=?, cook_time=?, " +
	"	total_time=?, category_id=?, keywords=?, yield=?, nutrition_id=?" +
	"WHERE id=?"

var updateDateModifiedStmt = "UPDATE " + schema.recipe.name + " " +
	"SET date_modified=? " +
	"WHERE id=?"

// DELETE statements
func deleteRecipeStmts() []string {
	return []string{
		deleteAssocValues(schema.ingredient),
		deleteAssocValues(schema.instruction),
		deleteAssocValues(schema.tool),
		"DELETE FROM " + schema.recipe.name + " WHERE id=?",
	}
}

func deleteAssocValues(t table) string {
	return "DELETE FROM " + t.assocTable + " WHERE recipe_id=?"
}
