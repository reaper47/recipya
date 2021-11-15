package model

// RecipesInfoWrapper holds a RecipesInfo object.
type RecipesInfoWrapper struct {
	Info *RecipesInfo `json:"info"`
}

// RecipesInfo holds metadata on recipes, i.e. the number of recipes.
type RecipesInfo struct {
	Total            int64            `json:"total"`
	TotalPerCategory map[string]int64 `json:"totalPerCategory"`
}
