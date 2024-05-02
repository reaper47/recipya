package statements

// InsertAdditionalImageRecipe is the query to add an image to a recipe.
const InsertAdditionalImageRecipe = `
	INSERT INTO additional_images_recipe (recipe_id, image)
	VALUES (?, ?)`

// InsertAuthToken is the query to add an authentication token to the database.
const InsertAuthToken = `
	INSERT INTO auth_tokens (selector, hash_validator, user_id)
	VALUES (?, ?, ?)`

// InsertCategory is the query to add a category to the database.
const InsertCategory = `
	INSERT INTO categories (name)
	VALUES (?)
	ON CONFLICT DO UPDATE SET name = EXCLUDED.name
	RETURNING id`

// InsertCookbook is the query to add a cookbook to the database.
const InsertCookbook = `
	INSERT INTO cookbooks (title, image, user_id) 
	VALUES (?, ?, ?)
	RETURNING id`

// InsertCookbookRecipe is the query to add a recipe to a cookbook.
const InsertCookbookRecipe = `
	INSERT INTO cookbook_recipes (cookbook_id, recipe_id, order_index)
	VALUES (?, ?, (SELECT COUNT(*)
				   FROM cookbooks AS c
							INNER JOIN cookbook_recipes AS cr ON cr.cookbook_id = c.id
				   WHERE c.id = ?
					 AND c.user_id = ?))`

// InsertCuisine is the query to add a cuisine to the database
const InsertCuisine = `
	INSERT OR IGNORE INTO cuisines (name)
	VALUES (?)`

// InsertReport is the query to add a report without logs into the database
const InsertReport = `
	INSERT INTO reports (report_type, created_at, exec_time_ns, user_id) 
	VALUES (?, ?, ?, ?)
	RETURNING id`

// InsertReportLog is the query to add a log to a report.
const InsertReportLog = `
	INSERT INTO report_logs (report_id, title, success, error_reason) 
	VALUES (?, ?, ?, ?)`

// InsertIngredient is the query to add an ingredient.
const InsertIngredient = `
	INSERT INTO ingredients (name)
	VALUES (?)
	ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
	RETURNING id`

// InsertInstruction is the query to add an instruction.
const InsertInstruction = `
	INSERT INTO instructions (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
	RETURNING id`

// InsertKeyword is the query to add a keyword.
const InsertKeyword = `
	INSERT INTO keywords (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
	RETURNING id`

// InsertNutrition is the query to add a nutrition facts.
const InsertNutrition = `
	INSERT INTO nutrition (recipe_id, calories, total_carbohydrates, sugars, protein, total_fat, saturated_fat, unsaturated_fat, cholesterol, sodium, fiber, is_per_serving)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

// InsertRecipe is the query to add a recipe to the database.
const InsertRecipe = `
	INSERT INTO recipes (name, description, image, yield, url)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id`

// InsertRecipeCategory associates a recipe with a category.
const InsertRecipeCategory = `
	INSERT INTO category_recipe (category_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeCuisine associates a recipe with a category.
const InsertRecipeCuisine = `
	INSERT INTO cuisine_recipe (cuisine_id, recipe_id)
	VALUES (?, ?)`

const InsertRecipeImage = `
	INSERT INTO additional_images_recipe (recipe_id, image)
	VALUES (?, ?)`

// InsertRecipeIngredient is the query to associate a recipe with an ingredient.
const InsertRecipeIngredient = `
	INSERT INTO ingredient_recipe (ingredient_id, recipe_id, ingredient_order)
	VALUES (?, ?, ?)`

// InsertRecipeInstruction is the query to associate a recipe with an instruction.
const InsertRecipeInstruction = `
	INSERT INTO instruction_recipe (instruction_id, recipe_id, instruction_order)
	VALUES (?, ?, ?)`

// InsertRecipeKeyword is the query to associate a recipe with a keyword.
const InsertRecipeKeyword = `
	INSERT INTO keyword_recipe (keyword_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeShadow is the query to insert a recipe into the shadow table.
const InsertRecipeShadow = `
	INSERT OR REPLACE INTO shadow_last_inserted_recipe (row, id, name, description, source)
	VALUES (1, ?, ?, ?, ?)`

// InsertRecipeTime is the query to associate a recipe with a time.
const InsertRecipeTime = `
	INSERT INTO time_recipe (time_id, recipe_id)
	VALUES (?, ?)`

// InsertRecipeTool is the query to associate a recipe with a tool.
const InsertRecipeTool = `
	INSERT INTO tool_recipe (tool_id, recipe_id)
	VALUES (?, ?)`

// InsertShareLink is the query to add a recipe share link to the database.
const InsertShareLink = `
	INSERT INTO share_recipes (link, recipe_id, user_id)
	VALUES (?, ?, ?)
	ON CONFLICT (link, recipe_id) DO NOTHING`

// InsertShareLinkCookbook is the query to add a cookbook share link to the database.
const InsertShareLinkCookbook = `
	INSERT INTO share_cookbooks (link, cookbook_id, user_id)
	VALUES (?, ?, ?)
	ON CONFLICT (link, cookbook_id) DO NOTHING`

// InsertTimes is the query to add kitchen times.
const InsertTimes = `
	INSERT INTO times (prep_seconds, cook_seconds)
	VALUES (?, ?)
	ON CONFLICT (prep_seconds, cook_seconds) 
	    DO UPDATE 
		SET prep_seconds = EXCLUDED.prep_seconds,
		    cook_seconds = EXCLUDED.cook_seconds
	RETURNING id`

// InsertTool is the query to add a tool.
const InsertTool = `
	INSERT INTO tools (name)
	VALUES (?)
	ON CONFLICT (name)
		DO UPDATE SET name = EXCLUDED.name
	RETURNING id`

// InsertUser is the query to add a user to the database.
const InsertUser = `
	INSERT INTO users (email, hashed_password)
	VALUES (?, ?)
	RETURNING id`

// InsertUserCategory is the query to associate a category with a user.
const InsertUserCategory = `
	INSERT OR IGNORE INTO user_category (user_id, category_id)
	VALUES (?, ?)`

// InsertUserRecipe is the query to associate a recipe with a user.
const InsertUserRecipe = `
	INSERT INTO user_recipe (user_id, recipe_id)
	VALUES (?, ?)`
