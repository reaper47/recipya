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
	Ingredients  map[string]string
	Instructions map[string]string
	Keywords     map[string]string
	Tools        map[string]string

	InsertIngredientsStmt  string
	InsertInstructionsStmt string
	InsertKeywordsStmt     string
	InsertToolsStmt        string
}

func (m *nameParams) init(tables []tableData, offset int) {
	m.Ingredients = make(map[string]string)
	m.Instructions = make(map[string]string)
	m.Keywords = make(map[string]string)
	m.Tools = make(map[string]string)

	m.InsertIngredientsStmt, offset = insertIntoNameTableStmt(
		"ingredients",
		tables[0].Entries,
		offset,
		m.Ingredients,
	)

	m.InsertInstructionsStmt, offset = insertIntoNameTableStmt(
		"instructions",
		tables[1].Entries,
		offset,
		m.Instructions,
	)

	m.InsertKeywordsStmt, offset = insertIntoNameTableStmt(
		"keywords",
		tables[2].Entries,
		offset,
		m.Keywords,
	)

	m.InsertToolsStmt, _ = insertIntoNameTableStmt(
		"tools",
		tables[3].Entries,
		offset,
		m.Tools,
	)
}

func (m *nameParams) insertStmts(tables []tableData, isInsRecipeDefined bool) string {
	return m.InsertIngredientsStmt + "" +
		insertIntoAssocTableStmt(tables[0], "ins_ingredients", m.Ingredients, isInsRecipeDefined) + "" +
		m.InsertInstructionsStmt + "" +
		insertIntoAssocTableStmt(tables[1], "ins_instructions", m.Instructions, isInsRecipeDefined) + "" +
		m.InsertKeywordsStmt + "" +
		insertIntoAssocTableStmt(tables[2], "ins_keywords", m.Keywords, isInsRecipeDefined) + "" +
		m.InsertToolsStmt + "" +
		insertIntoAssocTableStmt(tables[3], "ins_tools", m.Tools, isInsRecipeDefined)

}

// GetAllRecipes gets all of the recipes in the database.
func (m *postgresDBRepo) GetAllRecipes(page int) ([]models.Recipe, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	rows, err := m.Pool.Query(ctx, getRecipesPagination(page))
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

func (m *postgresDBRepo) GetRecipesCount() (int, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var count int
	err := m.Pool.QueryRow(ctx, recipesCountStmt).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
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

	tables := getTables(r)
	var recipeID int64
	err = tx.QueryRow(ctx, insertRecipeStmt(tables), r.ToArgs(false)...).Scan(&recipeID)
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

// UpdateRecipes update the recipe in the database.
func (m *postgresDBRepo) UpdateRecipe(r models.Recipe) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = m.Pool.Exec(ctx, deleteAssocTableEntries, r.ID)
	if err != nil {
		return err
	}

	tables := getTables(r)
	_, err = m.Pool.Exec(ctx, updateRecipeStmt(tables), r.ToArgs(true)...)
	if err != nil {
		return err
	}

	tables = append(tables, tableData{Table: "categories"})
	for _, td := range tables {
		_, err = m.Pool.Exec(ctx, resetIDStmt(td.Table))
		if err != nil {
			return err
		}
	}
	return nil
}

func getTables(r models.Recipe) []tableData {
	return []tableData{
		{Table: "ingredients", AssocTable: "ingredient_recipe", Entries: r.Ingredients},
		{Table: "instructions", AssocTable: "instruction_recipe", Entries: r.Instructions},
		{Table: "keywords", AssocTable: "keyword_recipe", Entries: r.Keywords},
		{Table: "tools", AssocTable: "tool_recipe", Entries: r.Tools},
	}
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
