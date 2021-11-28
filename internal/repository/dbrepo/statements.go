package db

import (
	"strconv"
	"strings"
)

var insertRecipeStmt = `
	WITH nutrition AS (
		INSERT INTO nutrition (
			calories, total_carbohydrates, sugars, protein,
			total_fat, saturated_fat, cholesterol, sodium, fiber
		) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	), all_categories (name) as (
		values ($10)
	), inserted_categories (id, name) as (
		INSERT INTO categories (name)
		SELECT name from all_categories
		ON CONFLICT DO NOTHING 
		RETURNING id, name 
	), times_to_insert (prep, cook) AS (
		VALUES ($11::interval,$12::interval)
	), inserted_times (id, prep, cook) AS (
		INSERT INTO times (prep, cook)
		SELECT prep, cook
		FROM times_to_insert
		ON CONFLICT (prep, cook) DO NOTHING
		RETURNING id, prep, cook
	)
	INSERT INTO recipes (
		name, description, url, image, yield, 
		category_id, times_id, nutrition_id
	)
	VALUES (
		$13,$14,$15,$16,$17,
		(
			SELECT id 
			FROM (
				SELECT c.id, c.name
				FROM categories c
				UNION ALL
				SELECT id, name
				FROM inserted_categories
			) AS id 
			where name=$18
		),
		(
			SELECT id 
			FROM (
				SELECT t.id, t.prep, t.cook FROM times t
				UNION ALL
				SELECT id, prep, cook FROM inserted_times
			) AS id
			WHERE prep=$19 AND cook=$20
		),
		(SELECT id FROM nutrition)
	)
	RETURNING id
`

func insertXsStmt(table string, xs []string) (string, []interface{}) {
	si := make([]interface{}, len(xs))
	var values, wheres string
	for i, s := range xs {
		values += "($" + strconv.Itoa(i+1) + ")"
		wheres += "name=$" + strconv.Itoa(i+1)
		if i < len(xs)-1 {
			values += ","
			wheres += " OR "
		}
		si[i] = s
	}

	sql := `
		WITH inserts AS(
			INSERT INTO ` + table + ` (name) 
			VALUES ` + values + ` 
			ON CONFLICT (name) DO NOTHING
			RETURNING id
		)
		SELECT * FROM inserts
		UNION
		SELECT id FROM ` + table + ` 
		WHERE ` + wheres

	return sql, si
}

func insertAssocStmt(assocTable string, recipeID int64, ids []int64) (string, []interface{}) {
	col := strings.SplitN(assocTable, "_", 2)[0] + "_id"
	values := ""
	si := make([]interface{}, len(ids))
	for i, id := range ids {
		values += "(" + strconv.FormatInt(recipeID, 10) + ",$" + strconv.Itoa(i+1) + ")"
		if i < len(ids)-1 {
			values += ","
		}
		si[i] = id
	}
	sql := "INSERT INTO " + assocTable + " (recipe_id," + col + ") VALUES " + values + " ON CONFLICT DO NOTHING"
	return sql, si
}

func resetIDStmt(table string) string {
	return "SELECT setval('" + table + "_id_seq', MAX(id)) FROM " + table
}
