package db

import (
	"time"

	"github.com/reaper47/recipya/internal/contexts"
	"github.com/reaper47/recipya/internal/models"
)

type tableData struct {
	Table      string
	AssocTable string
	Entries    []string
}

// GetAllRecipes gets all of the recipes in the database.
func (m *postgresDBRepo) GetAllRecipes() ([]models.Recipe, error) {
	recipes := []models.Recipe{}
	return recipes, nil
}

// InsertNewRecipe inserts a new recipe into the database.
//
// The CreatedAt and UpdatedAt timestamps are not inserted
// because the database will take care it.
func (m *postgresDBRepo) InsertNewRecipe(r models.Recipe) (int64, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	var recipeID int64
	err = tx.QueryRow(ctx, insertRecipeStmt, r.Nutrition.Calories,
		r.Nutrition.TotalCarbohydrates,
		r.Nutrition.Sugars,
		r.Nutrition.Protein,
		r.Nutrition.TotalFat,
		r.Nutrition.SaturatedFat,
		r.Nutrition.Cholesterol,
		r.Nutrition.Sodium,
		r.Nutrition.Fiber,
		r.Category,
		r.Times.Prep,
		r.Times.Cook,
		r.Name,
		r.Description,
		r.Url,
		r.Image,
		r.Yield,
		r.Category,
		r.Times.Prep,
		r.Times.Cook,
	).Scan(&recipeID)
	if err != nil {
		return -1, err
	}

	tables := []tableData{
		{Table: "ingredients", AssocTable: "ingredient_recipe", Entries: r.Ingredients},
		{Table: "instructions", AssocTable: "instruction_recipe", Entries: r.Instructions},
	}
	for _, entry := range tables {
		ids, err := func() ([]int64, error) {
			sql, si := insertXsStmt(entry.Table, entry.Entries)
			rows, err := tx.Query(ctx, sql, si...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var ids []int64
			for rows.Next() {
				var id int64
				err = rows.Scan(&id)
				if err != nil {
					return nil, err
				}
				ids = append(ids, id)
			}
			return ids, nil
		}()
		if err != nil {
			return -1, err
		}

		sql, si := insertAssocStmt(entry.AssocTable, recipeID, ids)
		_, err = tx.Exec(ctx, sql, si...)
		if err != nil {
			return -1, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}

	for _, td := range tables {
		_, err = m.Pool.Exec(ctx, resetIDStmt(td.Table))
		if err != nil {
			return -1, err
		}
	}
	return recipeID, nil
}
