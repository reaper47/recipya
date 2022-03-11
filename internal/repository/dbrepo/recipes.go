package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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

func (m *postgresDBRepo) Recipes(userID int64, page int) ([]models.Recipe, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var rows pgx.Rows
	var err error
	if page == -1 {
		rows, err = m.Pool.Query(ctx, getAllRecipesStmt, userID)
	} else {
		rows, err = m.Pool.Query(ctx, getRecipesStmt, userID, (page-1)*12)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (m *postgresDBRepo) Recipe(id int64) models.Recipe {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var (
		r models.Recipe
		n models.Nutrition
		t models.Times
	)

	err := m.Pool.QueryRow(ctx, getRecipeStmt, id).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.Url,
		&r.Image,
		&r.Yield,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.Category,
		&n.Calories,
		&n.TotalCarbohydrates,
		&n.Sugars,
		&n.Protein,
		&n.TotalFat,
		&n.SaturatedFat,
		&n.Cholesterol,
		&n.Sodium,
		&n.Fiber,
		&r.Ingredients,
		&r.Instructions,
		&r.Keywords,
		&r.Tools,
		&t.Prep,
		&t.Cook,
		&t.Total,
	)
	if err == pgx.ErrNoRows {
		return r
	}

	r.Nutrition = n
	r.Times = t
	return r
}

func (m *postgresDBRepo) RecipesCount() (int, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var count int
	err := m.Pool.QueryRow(ctx, recipesCountStmt).Scan(&count)
	if err == pgx.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return count, nil
}

func (m *postgresDBRepo) Categories(userID int64) []string {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var categories []string
	rows, err := m.Pool.Query(ctx, getCategoriesStmt, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c string
			err := rows.Scan(&c)
			if err != nil {
				continue
			}
			categories = append(categories, c)
		}
	}
	return categories
}

func (m *postgresDBRepo) InsertCategory(name string, userID int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	_, err := m.Pool.Exec(ctx, insertUserCategoryStmt, name, userID)
	return err
}

func (m *postgresDBRepo) InsertNewRecipe(r models.Recipe, userID int64) (int64, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	tables := getTables(r)
	args := []interface{}{userID}
	args = append(args, r.ToArgs(false)...)

	var recipeID int64
	if err = tx.QueryRow(ctx, insertRecipeStmt(tables), args...).Scan(&recipeID); err != nil {
		return -1, err
	}

	if err = tx.Commit(ctx); err != nil {
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
