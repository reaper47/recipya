package model

// Constants that point to specific FDC nutrient IDs.
const (
	FDCID_CALORIES      = 1008
	FDCID_CARBOHYDRATE  = 1005
	FDCID_FAT           = 1004
	FDCID_SATURATED_FAT = 1258
	FDCID_CHOLESTEROL   = 1253
	FDCID_PROTEIN       = 1003
	FDCID_SODIUM        = 1093
	FDCID_FIBER         = 1079
	FDCID_SUGAR         = 2000
)

// Foods is the wrapper for the foods JSON object.
type Foods struct {
	List []*Food `json:"foods"`
}

// Food stores relevant information on food items from the FDC API.
type Food struct {
	Id            int64            `json:"fdcId"`
	Description   string           `json:"description"`
	BrandName     string           `json:"brandName"`
	BrandOwner    string           `json:"brandOwner"`
	MarketCountry string           `json:"marketCountry"`
	Ingredients   string           `json:"ingredients"`
	Nutrients     []*FoodNutrients `json:"foodNutrients"`
	Score         float64          `json:"score"`
}

// FoodNutrients stores relevant information on the nutrition of
// food items from the FDC API.
type FoodNutrients struct {
	Id             int64   `json:"nutrientId"`
	Name           string  `json:"nutrientName"`
	Unit           string  `json:"unitName"`
	Amount         float64 `json:"amount"`
	Value          float64 `json:"value"`
	DerivationCode string  `json:"derivationCode"`
}

// IsOfInterest filters the FDC nutrient by ID.
func (f *FoodNutrients) IsOfInterest() bool {
	return f.Id == FDCID_CALORIES ||
		f.Id == FDCID_CARBOHYDRATE ||
		f.Id == FDCID_FAT ||
		f.Id == FDCID_SATURATED_FAT ||
		f.Id == FDCID_CHOLESTEROL ||
		f.Id == FDCID_PROTEIN ||
		f.Id == FDCID_SODIUM ||
		f.Id == FDCID_FIBER ||
		f.Id == FDCID_SUGAR
}
