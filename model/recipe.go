package model

// Recipes holds a list of Recipe objects.
type Recipes struct {
	Objects []*Recipe `json:"recipes"`
}

// Categories holds a list of categories.
type Categories struct {
	Objects []string `json:"categories"`
}

// Recipe holds the details of a recipe.
// It follows the Recipe schema standard.
type Recipe struct {
	ID                 int64         `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	Url                string        `json:"url"`
	Image              string        `json:"image"`
	PrepTime           string        `json:"prepTime"`
	CookTime           string        `json:"cookTime"`
	TotalTime          string        `json:"totalTime"`
	RecipeCategory     string        `json:"recipeCategory"`
	Keywords           string        `json:"keywords"`
	RecipeYield        int           `json:"recipeYield"`
	Tool               []string      `json:"tool"`
	RecipeIngredient   []string      `json:"recipeIngredient"`
	RecipeInstructions []string      `json:"recipeInstructions"`
	Nutrition          *NutritionSet `json:"nutrition"`
	DateModified       string        `json:"dateModified"`
	DateCreated        string        `json:"dateCreated"`
}

// NutritionSet holds nutritional details of a recipe.
type NutritionSet struct {
	FdcIDs       []int64 `json:"-"`
	Calories     string  `json:"calories"`
	Carbohydrate string  `json:"carbohydrateContent"`
	Fat          string  `json:"fatContent"`
	SaturatedFat string  `json:"saturatedFatContent"`
	Cholesterol  string  `json:"cholesterolContent"`
	Protein      string  `json:"proteinContent"`
	Sodium       string  `json:"sodiumContent"`
	Fiber        string  `json:"fiberContent"`
	Sugar        string  `json:"sugarContent"`
}

func (n *NutritionSet) IsEmpty() bool {
	return n.Calories == "" &&
		n.Carbohydrate == "" &&
		n.Fat == "" &&
		n.SaturatedFat == "" &&
		n.Cholesterol == "" &&
		n.Protein == "" &&
		n.Sodium == "" &&
		n.Fiber == "" &&
		n.Sugar == ""
}

func (r *Recipe) IsCreatedSameTime(other *Recipe) bool {
	return r.DateCreated == other.DateCreated
}

func (r *Recipe) IsModified(other *Recipe) bool {
	return r.DateModified != other.DateModified
}
