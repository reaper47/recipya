package statements

// InsertAuthToken is the query to add an authentication token to the database.
const InsertAuthToken = `
	INSERT INTO auth_tokens (selector, hash_validator, user_id)
	VALUES (?, ?, ?)`

// InsertCategory is the query to add a category to the database
const InsertCategory = `
	INSERT OR IGNORE INTO categories (name)
	VALUES (?)`

// InsertCuisine is the query to add a cuisine to the database
const InsertCuisine = `
	INSERT OR IGNORE INTO cuisines (name)
	VALUES (?)`

// InsertIngredient is the query to add an ingredient.
const InsertIngredient = `
	INSERT INTO ingredients (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = excluded.name`

// InsertInstruction is the query to add an instruction.
const InsertInstruction = `
	INSERT INTO instructions (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = excluded.name`

// InsertKeyword is the query to add a keyword.
const InsertKeyword = `
	INSERT INTO keywords (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = excluded.name`

// InsertNutrition is the query to add a nutrition fact.
const InsertNutrition = `
	INSERT INTO nutrition (recipe_id, calories, total_carbohydrates, sugars, protein, total_fat, saturated_fat, unsaturated_fat, cholesterol, sodium, fiber)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// InsertRecipe is the query to add a recipe to the database.
const InsertRecipe = `
	INSERT INTO recipes (name, description, image, yield, url)
	VALUES (?, ?, ?, ?, ?)`

// InsertRecipeCategory associates a recipe with a category.
const InsertRecipeCategory = `
	INSERT INTO category_recipe (category_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeCuisine associates a recipe with a category.
const InsertRecipeCuisine = `
	INSERT INTO cuisine_recipe (cuisine_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeIngredient is the query to associate a recipe with an ingredient.
const InsertRecipeIngredient = `
	INSERT INTO ingredient_recipe (ingredient_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeInstruction is the query to associate a recipe with an instruction.
const InsertRecipeInstruction = `
	INSERT INTO instruction_recipe (instruction_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeKeyword is the query to associate a recipe with a keyword.
const InsertRecipeKeyword = `
	INSERT INTO keyword_recipe (keyword_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeTime is the query to associate a recipe with a time.
const InsertRecipeTime = `
	INSERT INTO time_recipe (time_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeTool is the query to associate a recipe with a tool.
const InsertRecipeTool = `
	INSERT INTO tool_recipe (tool_id, recipe_id)
	VALUES (?, ?)`

// InsertTimes is the query to add kitchen times.
const InsertTimes = `
	INSERT INTO times (prep_seconds, cook_seconds)
	VALUES (?, ?)
	ON CONFLICT (prep_seconds, cook_seconds) DO UPDATE
		SET prep_seconds=0
	WHERE FALSE
	RETURNING id, prep_seconds, cook_seconds, total_seconds`

// InsertTool is the query to add a tool.
const InsertTool = `
	INSERT INTO tools (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = excluded.name`

// InsertUser is the query to add a user to the database.
const InsertUser = `
	INSERT INTO users (email, hashed_password)
	VALUES (?, ?)`

// InsertUserRecipe is the query to associate a recipe with a user.
const InsertUserRecipe = `
	INSERT INTO user_recipe (user_id, recipe_id)
	VALUES (?, ?)`
