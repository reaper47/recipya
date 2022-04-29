package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/reaper47/recipya/internal/contexts"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/nlp"
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

func (n *nameParams) init(tables []tableData, offset int) {
	n.Ingredients = make(map[string]string)
	n.Instructions = make(map[string]string)
	n.Keywords = make(map[string]string)
	n.Tools = make(map[string]string)

	n.InsertIngredientsStmt, offset = insertIntoNameTableStmt(
		"ingredients",
		tables[0].Entries,
		offset,
		n.Ingredients,
	)

	n.InsertInstructionsStmt, offset = insertIntoNameTableStmt(
		"instructions",
		tables[1].Entries,
		offset,
		n.Instructions,
	)

	n.InsertKeywordsStmt, offset = insertIntoNameTableStmt(
		"keywords",
		tables[2].Entries,
		offset,
		n.Keywords,
	)

	n.InsertToolsStmt, _ = insertIntoNameTableStmt(
		"tools",
		tables[3].Entries,
		offset,
		n.Tools,
	)
}

func (n *nameParams) insertStmts(tables []tableData, isInsRecipeDefined bool) string {
	return n.InsertIngredientsStmt + "" +
		insertIntoAssocTableStmt(tables[0], "ins_ingredients", n.Ingredients, isInsRecipeDefined) + "" +
		n.InsertInstructionsStmt + "" +
		insertIntoAssocTableStmt(tables[1], "ins_instructions", n.Instructions, isInsRecipeDefined) + "" +
		n.InsertKeywordsStmt + "" +
		insertIntoAssocTableStmt(tables[2], "ins_keywords", n.Keywords, isInsRecipeDefined) + "" +
		n.InsertToolsStmt + "" +
		insertIntoAssocTableStmt(tables[3], "ins_tools", n.Tools, isInsRecipeDefined)

}

func (p *postgresDBRepo) Recipes(userID int64, page int) (xr []models.Recipe, err error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var rows pgx.Rows
	if page == -1 {
		rows, err = p.Pool.Query(ctx, getAllRecipesStmt, userID)
	} else {
		rows, err = p.Pool.Query(ctx, getRecipesStmt, userID, (page-1)*12)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			URL:         url,
			Yield:       yield,
			Category:    category,
			Times: models.Times{
				Prep:  prep,
				Cook:  cook,
				Total: total,
			},
			Ingredients: models.Ingredients{Values: ingredients},
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

func (p *postgresDBRepo) RecipeForUser(id, userID int64) (models.Recipe, error) {
	if !p.RecipeBelongsToUser(id, userID) {
		return models.Recipe{}, errors.New("recipe does not belong to user")
	}
	return p.Recipe(id), nil
}

func (p *postgresDBRepo) RecipeBelongsToUser(id, userID int64) bool {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var (
		u int64
		r int64
	)
	err := p.Pool.QueryRow(ctx, getUserRecipe, userID, id).Scan(&u, &r)
	return err != pgx.ErrNoRows
}

func (p *postgresDBRepo) Recipe(id int64) models.Recipe {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var (
		r models.Recipe
		n models.Nutrition
		t models.Times
	)

	err := p.Pool.QueryRow(ctx, getRecipeStmt, id).Scan(
		&r.ID,
		&r.Name,
		&r.Description,
		&r.URL,
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
		&r.Ingredients.Values,
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
	r.Normalize()
	return r
}

func (p *postgresDBRepo) RecipesCount() (int, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	var count int
	err := p.Pool.QueryRow(ctx, recipesCountStmt).Scan(&count)
	if err == pgx.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return -1, err
	}
	return count, nil
}

func (p *postgresDBRepo) Categories(userID int64) (xs []string) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	rows, err := p.Pool.Query(ctx, getCategoriesStmt, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c string
			err := rows.Scan(&c)
			if err != nil {
				continue
			}
			xs = append(xs, c)
		}
	}
	return xs
}

func (p *postgresDBRepo) EditCategory(oldCategory, newCategory string, userID int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	_, err := p.Pool.Exec(ctx, updateCategoryStmt, oldCategory, newCategory, userID)
	return err
}

func (p *postgresDBRepo) InsertCategory(name string, userID int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	_, err := p.Pool.Exec(ctx, insertUserCategoryStmt, name, userID)
	return err
}

func (p *postgresDBRepo) DeleteCategory(name string, userID int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	_, err := p.Pool.Exec(ctx, deleteUserCategoryStmt, userID, name)
	if err != nil {
		return fmt.Errorf("delete category - error deleting user category: %s", err)
	}

	_, err = p.Pool.Exec(ctx, deleteRecipeCategoryStmt, userID, name)
	if err != nil {
		return fmt.Errorf("delete category - error deleting recipe category: %s", err)
	}
	return nil
}

func (p *postgresDBRepo) InsertNewRecipe(r models.Recipe, userID int64) (int64, error) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	var isExists bool
	err = tx.QueryRow(ctx, isRecipeExistsForUserStmt, userID, r.Name, r.Description, r.URL, r.Yield).Scan(&isExists)
	if err != nil {
		return -1, err
	}

	if isExists {
		return -1, fmt.Errorf("recipe '%s' exists for user %d", r.Name, userID)
	}

	r.Instructions = nlp.CapitalizeParagraphs(r.Instructions)

	tables := getTables(r)
	args := []interface{}{userID}
	args = append(args, r.ToArgs(false)...)

	var recipeID int64
	err = tx.QueryRow(ctx, insertRecipeStmt(tables), args...).Scan(&recipeID)
	if err != nil {
		return -1, err
	}

	if err = tx.Commit(ctx); err != nil {
		return -1, err
	}

	for _, td := range tables {
		_, err = p.Pool.Exec(ctx, resetIDStmt(td.Table))
		if err != nil {
			return -1, err
		}
	}
	return recipeID, nil
}

func (p *postgresDBRepo) UpdateRecipe(r models.Recipe) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = p.Pool.Exec(ctx, deleteAssocTableEntries, r.ID)
	if err != nil {
		return err
	}

	r.Instructions = nlp.CapitalizeParagraphs(r.Instructions)

	tables := getTables(r)
	_, err = p.Pool.Exec(ctx, updateRecipeStmt(tables), r.ToArgs(true)...)
	if err != nil {
		return err
	}

	tables = append(tables, tableData{Table: "categories"})
	for _, td := range tables {
		_, err = p.Pool.Exec(ctx, resetIDStmt(td.Table))
		if err != nil {
			return err
		}
	}
	return nil
}

func getTables(r models.Recipe) []tableData {
	return []tableData{
		{Table: "ingredients", AssocTable: "ingredient_recipe", Entries: r.Ingredients.Values},
		{Table: "instructions", AssocTable: "instruction_recipe", Entries: r.Instructions},
		{Table: "keywords", AssocTable: "keyword_recipe", Entries: r.Keywords},
		{Table: "tools", AssocTable: "tool_recipe", Entries: r.Tools},
	}
}

func (p *postgresDBRepo) DeleteRecipe(id int64) error {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	cmd, err := p.Pool.Exec(ctx, deleteRecipeStmt, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("recipe not found")
	}
	return nil
}

func (p *postgresDBRepo) Images() (xs []string) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	rows, err := p.Pool.Query(ctx, getDistinctImagesStmt)
	if err != nil {
		return xs
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		_ = rows.Scan(&s)
		xs = append(xs, s)
	}
	return xs
}

func (p *postgresDBRepo) Websites() (xs []models.Website) {
	ctx, cancel := contexts.Timeout(3 * time.Second)
	defer cancel()

	rows, err := p.Pool.Query(ctx, getWebsitesStmt)
	if err != nil {
		return xs
	}
	defer rows.Close()

	for rows.Next() {
		var name, url string
		_ = rows.Scan(&name, &url)
		xs = append(xs, models.Website{
			Name: name,
			URL:  url,
		})
	}
	return xs
}
