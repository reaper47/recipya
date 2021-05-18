package model

type Recipes struct {
	Objects []*Recipe `json:"recipes"`
}

type Recipe struct {
	ID				   int64
	Name               string     
	Description        string       
	Url                string       
	Image              string       
	PrepTime           string       
	CookTime           string       
	TotalTime          string       
	RecipeCategory     string       
	Keywords           string       
	RecipeYield        int          
	Tool               []string     
	RecipeIngredient   []string     
	RecipeInstructions []string     
	Nutrition          *NutritionSet 
	DateModified       string       
	DateCreated        string       
}

type NutritionSet struct {
	Calories     string 
	Carbohydrate string `json:"carbohydrateContent"`
	Fat          string `json:"fatContent"`
	SaturatedFat string `json:"saturatedFatContent"`
	Cholesterol  string `json:"cholesterolContent"`
	Protein      string `json:"proteinContent"`
	Sodium       string `json:"sodiumContent"`
	Fiber        string `json:"fiberContent"`
	Sugar        string `json:"sugarContent"`
}
