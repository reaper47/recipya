package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/contexts"
	"github.com/reaper47/recipya/internal/models"
)

type tableData struct {
	Table      string
	AssocTable string
	Entries    []string
}

type nameParams struct {
	Ingredients, Instructions, Keywords, Tools map[string]string
}

func (m *nameParams) init() {
	m.Ingredients = make(map[string]string)
	m.Instructions = make(map[string]string)
	m.Keywords = make(map[string]string)
	m.Tools = make(map[string]string)
}

// GetAllRecipes gets all of the recipes in the database.
func (m *postgresDBRepo) GetAllRecipes() ([]models.Recipe, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	rows, err := m.Pool.Query(ctx, getRecipes(false))
	if err != nil {
		return nil, err
	}

	var xr []models.Recipe
	for rows.Next() {
		var (
			recipeID                                           int64
			name, description, url, category                   string
			image                                              uuid.UUID
			yield                                              int16
			prep, cook, total                                  time.Duration
			calories, totalCarbohydrates, sugars, protein      string
			totalFat, saturatedFat, cholesterol, sodium, fiber string
			ingredients, instructions, keywords, tools         []string
			createdAt, updatedAt                               time.Time
		)

		err := rows.Scan(
			&recipeID,
			&name,
			&description,
			&url,
			&image,
			&yield,
			&createdAt,
			&updatedAt,
			&category,
			&calories,
			&totalCarbohydrates,
			&sugars,
			&protein,
			&totalFat,
			&saturatedFat,
			&cholesterol,
			&sodium,
			&fiber,
			&ingredients,
			&instructions,
			&keywords,
			&tools,
			&prep,
			&cook,
			&total,
		)
		if err != nil {
			return nil, err
		}

		xr = append(xr, models.Recipe{
			ID:          recipeID,
			Name:        name,
			Description: description,
			Image:       image,
			Url:         url,
			Yield:       yield,
			Category:    category,
			Times: models.Times{
				Prep:  prep,
				Cook:  cook,
				Total: total,
			},
			Ingredients: ingredients,
			Nutrition: models.Nutrition{
				Calories:           calories,
				TotalCarbohydrates: totalCarbohydrates,
				Sugars:             sugars,
				Protein:            protein,
				TotalFat:           totalFat,
				SaturatedFat:       saturatedFat,
				Cholesterol:        cholesterol,
				Sodium:             sodium,
				Fiber:              fiber,
			},
			Instructions: instructions,
			Keywords:     keywords,
			Tools:        tools,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		})
	}
	return xr, nil
}

// GetRecipe gets a recipe in the database.
func (m *postgresDBRepo) GetRecipe(id int64) (models.Recipe, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var (
		recipeID                                           int64
		name, description, url, category                   string
		image                                              uuid.UUID
		yield                                              int16
		prep, cook, total                                  time.Duration
		calories, totalCarbohydrates, sugars, protein      string
		totalFat, saturatedFat, cholesterol, sodium, fiber string
		ingredients, instructions, keywords, tools         []string
		createdAt, updatedAt                               time.Time
	)

	err := m.Pool.QueryRow(ctx, getRecipes(true), id).Scan(
		&recipeID,
		&name,
		&description,
		&url,
		&image,
		&yield,
		&createdAt,
		&updatedAt,
		&category,
		&calories,
		&totalCarbohydrates,
		&sugars,
		&protein,
		&totalFat,
		&saturatedFat,
		&cholesterol,
		&sodium,
		&fiber,
		&ingredients,
		&instructions,
		&keywords,
		&tools,
		&prep,
		&cook,
		&total,
	)
	if err != nil {
		return models.Recipe{}, err
	}

	return models.Recipe{
		ID:          recipeID,
		Name:        name,
		Description: description,
		Image:       image,
		Url:         url,
		Yield:       yield,
		Category:    category,
		Times: models.Times{
			Prep:  prep,
			Cook:  cook,
			Total: total,
		},
		Ingredients: ingredients,
		Nutrition: models.Nutrition{
			Calories:           calories,
			TotalCarbohydrates: totalCarbohydrates,
			Sugars:             sugars,
			Protein:            protein,
			TotalFat:           totalFat,
			SaturatedFat:       saturatedFat,
			Cholesterol:        cholesterol,
			Sodium:             sodium,
			Fiber:              fiber,
		},
		Instructions: instructions,
		Keywords:     keywords,
		Tools:        tools,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
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

	tables := []tableData{
		{Table: "ingredients", AssocTable: "ingredient_recipe", Entries: r.Ingredients},
		{Table: "instructions", AssocTable: "instruction_recipe", Entries: r.Instructions},
		{Table: "keywords", AssocTable: "keyword_recipe", Entries: r.Keywords},
		{Table: "tools", AssocTable: "tool_recipe", Entries: r.Tools},
	}

	var recipeID int64
	stmt := insertRecipeStmt(tables)
	err = tx.QueryRow(ctx, stmt, r.ToArgs()...).Scan(&recipeID)
	if err != nil {
		return -1, err
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

// DeleteRecipe deletes the recipe with the passed id from the database.
func (m *postgresDBRepo) DeleteRecipe(id int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	cmd, err := m.Pool.Exec(ctx, deleteRecipeStmt, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("recipe not found")
	}
	return nil
}
