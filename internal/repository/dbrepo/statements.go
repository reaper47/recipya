package db

import (
	"strconv"
	"strings"
)

// SELECT
func getRecipes(onlyOne bool) string {
	id := "r.id"
	if onlyOne {
		id = "$1"
	}

	stmt := `SELECT 
		r.id,
		r.name,
		description,
		url,
		image,
		yield,
		created_at,
		updated_at,
		c.name AS category,
		n.calories,
		n.total_carbohydrates,
		n.sugars,
		n.protein,
		n.total_fat,
		n.saturated_fat,
		n.cholesterol,
		n.sodium,
		n.fiber,
		ARRAY(
			SELECT name
			FROM ingredients i
			JOIN ingredient_recipe ir ON ir.ingredient_id = i.id
			WHERE ir.recipe_id = ` + id + `
		) AS ingredients,
		ARRAY(
			SELECT name
			FROM instructions i2
			JOIN instruction_recipe ir2 ON ir2.instruction_id = i2.id
			WHERE ir2.recipe_id = ` + id + `
		) AS instructions,
		ARRAY(
			SELECT name
			FROM keywords k
			JOIN keyword_recipe kr ON kr.keyword_id = k.id
			WHERE kr.recipe_id = ` + id + `
		) AS keywords,
		ARRAY(
			SELECT name
			FROM tools t
			JOIN tool_recipe tr ON tr.tool_id = t.id
			WHERE tr.recipe_id = ` + id + `
		) AS tools,
		t2.prep,
		t2.cook,
		t2.total
	FROM recipes r
	JOIN category_recipe cr ON cr.recipe_id = r.id
	JOIN categories c ON c.id = cr.category_id
	JOIN nutrition n ON n.recipe_id = r.id
	JOIN time_recipe tr2 ON tr2.recipe_id = r.id
	JOIN times t2 ON t2.id = tr2.time_id`

	if onlyOne {
		stmt += " WHERE r.id = $1"
	}
	return stmt
}

func resetIDStmt(table string) string {
	return "SELECT setval('" + table + "_id_seq', MAX(id)) FROM " + table
}

// INSERT
func insertRecipeStmt(tables []tableData) string {
	var params nameParams
	params.init(tables, 18)

	return `
		WITH ins_recipe AS (
			INSERT  INTO recipes (name, description, image, url, yield)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		), ins_category AS (
			INSERT INTO categories (name)
			VALUES ($6)
			ON CONFLICT ON CONSTRAINT categories_name_key DO UPDATE
			SET name=NULL
			WHERE FALSE
			RETURNING id, name
		), ins_category_id AS (
			INSERT INTO category_recipe (recipe_id, category_id)
			VALUES (
				(
					SELECT id 
					FROM ins_recipe
				),
				(
					SELECT id FROM ins_category
					UNION ALL
					SELECT id
					FROM categories
					WHERE name=$6
				)
			)
		),  ins_nutrition AS (
			INSERT INTO nutrition (
				recipe_id, calories, total_carbohydrates, sugars,
				protein, total_fat, saturated_fat, cholesterol, sodium, fiber
			)
			VALUES ((SELECT id FROM ins_recipe),$7,$8,$9,$10,$11,$12,$13,$14,$15)
			RETURNING id
		),  ins_times AS (
			INSERT INTO times (prep, cook)
			VALUES ($16::interval, $17::interval)
			ON CONFLICT ON CONSTRAINT times_prep_cook_key DO UPDATE
			SET prep=NULL
			WHERE FALSE
			RETURNING id, prep, cook, total
		), ins_time_recipe AS (
			INSERT INTO time_recipe (time_id, recipe_id)
			VALUES (
				(
					SELECT id FROM ins_times WHERE prep=$16::interval and cook=$17::interval
					UNION ALL
					SELECT id FROM times WHERE prep=$16::interval and cook=$17::interval
				),
				(
					SELECT id
					FROM ins_recipe
				)
			)
		)` + params.insertStmts(tables, true) + `
	SELECT id FROM ins_recipe`
}

func insertIntoNameTableStmt(
	name string,
	values []string,
	offset int,
	params map[string]string,
) (string, int) {
	if len(values) == 0 {
		return "", offset
	}

	var stmt = ", ins_" + name + ` AS (
		INSERT INTO ` + name + " (name) VALUES "

	for i, v := range values {
		param := "$" + strconv.Itoa(offset)
		stmt += "(" + param + ")"
		if i < len(values)-1 {
			stmt += ","
		}
		params[v] = param
		offset++
	}

	stmt += `
		ON CONFLICT ON CONSTRAINT ` + name + `_name_key DO UPDATE
		SET name=NULL
		WHERE false
		RETURNING id, name
	)`

	return stmt, offset
}

func insertIntoAssocTableStmt(td tableData, from string, params map[string]string, isInsRecipeDefiend bool) string {
	if len(td.Entries) == 0 {
		return ""
	}

	col := strings.SplitN(td.AssocTable, "_", 2)[0]
	tname := "ins_" + col + "_recipe"

	var stmt = "," + tname + ` AS (
		INSERT INTO ` + td.AssocTable + " (" + col + `_id, recipe_id) VALUES `

	var recipeID string
	if isInsRecipeDefiend {
		recipeID = "(SELECT id FROM ins_recipe)"
	} else {
		recipeID = "$1"
	}

	for i, v := range td.Entries {
		where := "WHERE name=" + params[v]
		stmt += `
		(
			(
				SELECT id FROM ` + from + " " + where + `
				UNION ALL
				SELECT id FROM ` + td.Table + " " + where + `
			),
			` + recipeID + `
		)`
		if i < len(td.Entries)-1 {
			stmt += ","
		}

	}
	return stmt + ")"
}

// UPDATE
func updateRecipeStmt(tables []tableData) string {
	var params nameParams
	params.init(tables, 19)

	return `
		WITH ins_category AS (
			INSERT INTO categories (name)
			VALUES ($7)
			ON CONFLICT ON CONSTRAINT categories_name_key DO UPDATE
			SET name=NULL
			WHERE FALSE
			RETURNING id, name
		), ins_category_id AS (
			UPDATE category_recipe 
			SET 
				category_id = (
					SELECT id FROM ins_category
					UNION ALL
					SELECT id FROM categories WHERE name=$7
				)
			WHERE recipe_id = $1
		), ins_nutrition AS (
			UPDATE nutrition 
			SET 
				calories = $8, 
				total_carbohydrates = $9, 
				sugars = $10,
				protein = $11,
				total_fat = $12,
				saturated_fat = $13,
				cholesterol = $14,
				sodium = $15,
				fiber = $16
			WHERE recipe_id = $1
		), ins_times AS (
			INSERT INTO times (prep, cook)
			VALUES ($17::interval, $18::interval)
			ON CONFLICT ON CONSTRAINT times_prep_cook_key DO UPDATE
			SET prep = NULL
			WHERE FALSE
			RETURNING id, prep, cook, total
		), ins_time_recipe AS (
			UPDATE time_recipe 
			SET 
				time_id = (
					SELECT id FROM ins_times WHERE prep=$17::interval and cook=$18::interval
					UNION ALL
					SELECT id FROM times WHERE prep=$17::interval and cook=$18::interval
				)
			WHERE recipe_id = $1
		)` + params.insertStmts(tables, false) + `

		UPDATE recipes r SET 
			name = $2,
			description = $3,
			image = CASE 
					WHEN 
						$4 != uuid_nil() AND
						$4 != image
					THEN $4
				ELSE 
					image
				END,
			url = $5,
			yield = $6
		WHERE id = $1`
}

// DELETE
const deleteRecipeStmt = "DELETE FROM recipes WHERE id = $1"

const deleteAssocTableEntries = `
	WITH del_ingredients AS (
		DELETE FROM ingredient_recipe
		WHERE recipe_id = $1
	), del_instructions AS (
		DELETE FROM instruction_recipe
		WHERE recipe_id = $1
	), del_tools AS (
		DELETE FROM tool_recipe
		WHERE recipe_id = $1
	)
	
	DELETE FROM keyword_recipe
	WHERE recipe_id = $1
`
