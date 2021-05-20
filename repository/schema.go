package repository

type tables struct {
	categories          table
	ingredients         table
	instructions        table
	nutrition           table
	recipes             table
	recipesIngredients  table
	recipesInstructions table
	recipesTools        table
	tools               table
}

type table struct {
	name       string
	assocTable string
	cols       map[string]string
}

var (
	foreignKeyExt     = "ON DELETE CASCADE ON UPDATE NO ACTION"
	recipesForeignKey = "FOREIGN KEY(recipes_id) REFERENCES recipes(id) " + foreignKeyExt
)

var schema = tables{
	categories: table{
		name:       "categories",
		assocTable: "recipes_categories",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	ingredients: table{
		name:       "ingredients",
		assocTable: "recipes_ingredients",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	instructions: table{
		name:       "instructions",
		assocTable: "recipes_instructions",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	nutrition: table{
		name: "nutrition",
		cols: map[string]string{
			"id":            "INTEGER PRIMARY KEY",
			"calories":      "TEXT",
			"carbohydrate":  "TEXT",
			"fat":           "TEXT",
			"saturated_fat": "TEXT",
			"cholesterol":   "TEXT",
			"protein":       "TEXT",
			"sodium":        "TEXT",
			"fiber":         "TEXT",
			"sugar":         "TEXT",
		},
	},
	recipes: table{
		name: "recipes",
		cols: map[string]string{
			"id":            "INTEGER PRIMARY KEY",
			"name":          "TEXT UNIQUE",
			"description":   "TEXT",
			"categories_id": "INTEGER",
			"nutrition_id":  "INTEGER",
			"url":           "TEXT",
			"image":         "TEXT",
			"prepTime":      "TEXT",
			"cookTime":      "TEXT",
			"totalTime":     "TEXT",
			"keywords":      "TEXT",
			"recipeYield":   "INTEGER",
			"dateModified":  "TEXT",
			"dateCreated":   "TEXT",
			"!foreignkey": "FOREIGN KEY (categories_id) REFERENCES categories(id) ON DELETE NO ACTION," +
				"FOREIGN KEY (nutrition_id) REFERENCES nutrition(id) ON DELETE NO ACTION",
		},
	},
	recipesIngredients: table{
		name: "recipes_ingredients",
		cols: map[string]string{
			"recipes_id":     "INTEGER",
			"ingredients_id": "INTEGER",
			"!foreignkey":    recipesForeignKey + ", " + "FOREIGN KEY(ingredients_id) REFERENCES ingredients(id) " + foreignKeyExt,
		},
	},
	recipesInstructions: table{
		name: "recipes_instructions",
		cols: map[string]string{
			"recipes_id":      "INTEGER",
			"instructions_id": "INTEGER",
			"!foreignkey":     recipesForeignKey + ", " + "FOREIGN KEY(instructions_id) REFERENCES instructions(id) " + foreignKeyExt,
		},
	},
	recipesTools: table{
		name: "recipes_tools",
		cols: map[string]string{
			"recipes_id":  "INTEGER",
			"tools_id":    "INTEGER",
			"!foreignkey": recipesForeignKey + ", " + "FOREIGN KEY(tools_id) REFERENCES tools(id) " + foreignKeyExt,
		},
	},
	tools: table{
		name:       "tools",
		assocTable: "recipes_tools",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
}

var allTables = []table{
	schema.categories,
	schema.ingredients,
	schema.instructions,
	schema.nutrition,
	schema.recipes,
	schema.recipesIngredients,
	schema.recipesInstructions,
	schema.recipesTools,
	schema.tools,
}
