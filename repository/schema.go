package repository

type tables struct {
	category          table
	ingredient        table
	instruction       table
	nutrition         table
	recipe            table
	recipeIngredient  table
	recipeInstruction table
	recipeTool        table
	tool              table
	website           table
}

type table struct {
	name       string
	assocTable string
	cols       map[string]string
}

const (
	foreignKeyExt    = "ON DELETE CASCADE ON UPDATE NO ACTION"
	recipeForeignKey = "FOREIGN KEY(recipe_id) REFERENCES recipe(id) " + foreignKeyExt
)

var schema = tables{
	category: table{
		name:       "category",
		assocTable: "recipe_category",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	ingredient: table{
		name:       "ingredient",
		assocTable: "recipe_ingredient",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	instruction: table{
		name:       "instruction",
		assocTable: "recipe_instruction",
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
	recipe: table{
		name: "recipe",
		cols: map[string]string{
			"id":            "INTEGER PRIMARY KEY",
			"name":          "TEXT UNIQUE",
			"description":   "TEXT",
			"category_id":   "INTEGER",
			"nutrition_id":  "INTEGER",
			"url":           "TEXT",
			"image":         "TEXT",
			"prep_time":     "TEXT",
			"cook_time":     "TEXT",
			"total_time":    "TEXT",
			"keywords":      "TEXT",
			"yield":         "INTEGER",
			"date_modified": "TEXT",
			"date_created":  "TEXT",
			"!foreignkey": "FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE NO ACTION," +
				"FOREIGN KEY (nutrition_id) REFERENCES nutrition(id) ON DELETE NO ACTION",
		},
	},
	recipeIngredient: table{
		name:       "recipe_ingredient",
		assocTable: "ingredient",
		cols: map[string]string{
			"recipe_id":     "INTEGER",
			"ingredient_id": "INTEGER",
			"!foreignkey":   recipeForeignKey + ", " + "FOREIGN KEY(ingredient_id) REFERENCES ingredient(id) " + foreignKeyExt,
		},
	},
	recipeInstruction: table{
		name:       "recipe_instruction",
		assocTable: "instruction",
		cols: map[string]string{
			"recipe_id":      "INTEGER",
			"instruction_id": "INTEGER",
			"!foreignkey":    recipeForeignKey + ", " + "FOREIGN KEY(instruction_id) REFERENCES instruction(id) " + foreignKeyExt,
		},
	},
	recipeTool: table{
		name:       "recipe_tool",
		assocTable: "tool",
		cols: map[string]string{
			"recipe_id":   "INTEGER",
			"tool_id":     "INTEGER",
			"!foreignkey": recipeForeignKey + ", " + "FOREIGN KEY(tool_id) REFERENCES tool(id) " + foreignKeyExt,
		},
	},
	tool: table{
		name:       "tool",
		assocTable: "recipe_tool",
		cols: map[string]string{
			"id":   "INTEGER PRIMARY KEY",
			"name": "TEXT NOT NULL UNIQUE",
		},
	},
	website: table{
		name: "website",
		cols: map[string]string{
			"id":  "INTEGER PRIMARY KEY",
			"url": "TEXT NOT NULL UNIQUE",
		},
	},
}

var allTables = []table{
	schema.category,
	schema.ingredient,
	schema.instruction,
	schema.nutrition,
	schema.recipe,
	schema.recipeIngredient,
	schema.recipeInstruction,
	schema.recipeTool,
	schema.tool,
	schema.website,
}
