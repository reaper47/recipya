package repository

import (
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
var selectRecipeByNameStmt = "SELECT " +
	"	id, name, description, url, image, prep_time, cook_time, total_time, " +
	"   (SELECT name FROM category WHERE id=category_id) AS category, " +
	"	keywords, yield, nutrition_id, date_modified, date_created " +
	"FROM " + schema.recipe.name + " " +
	"WHERE name=?"

var selectCategoryIdStmt = "SELECT id " +
	"FROM " + schema.category.name + " " +
	"WHERE name=?"

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

// INSERT statements
var insertRecipeStmt = "INSERT INTO " + schema.recipe.name + " (" +
	"name, description, url, image, prep_time, cook_time, " +
	"total_time, category_id, keywords, yield, " +
	"nutrition_id, date_modified, date_created" +
	") VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"

var insertNutritionStmt = "INSERT INTO " + schema.nutrition.name + " (" +
	"calories, carbohydrate, fat, saturated_fat, " +
	"cholesterol, protein, sodium, fiber, sugar" +
	") VALUES (?,?,?,?,?,?,?,?,?)"

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
func deleteAssocValues(t table) string {
	return "DELETE FROM " + t.assocTable + " WHERE recipe_id=?"
}
