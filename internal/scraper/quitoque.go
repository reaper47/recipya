package scraper

import (
	"encoding/json"
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type quitoque struct {
	Data struct {
		Recipe struct {
			Name             string `json:"name"`
			ShortDescription string `json:"shortDescription"`
			Image            string `json:"image"`
			Pools            []struct {
				NbPerson     int `json:"nbPerson"`
				CookingModes []struct {
					Name        string `json:"name"`
					CookingTime int    `json:"cookingTime"`
					WaitingTime int    `json:"waitingTime"`
					Stacks      struct {
						Tools []struct {
							ID              string `json:"id"`
							Name            string `json:"name"`
							Position        int    `json:"position"`
							QuantityToGram  any    `json:"quantityToGram"`
							QuantityToOrder any    `json:"quantityToOrder"`
							LiteralQuantity string `json:"literalQuantity"`
							Quantity        int    `json:"quantity"`
							Product         struct {
								ID                   string `json:"id"`
								Slug                 any    `json:"slug"`
								OriginalID           string `json:"originalId"`
								WeekDetailID         any    `json:"weekDetailId"`
								Nutriscore           any    `json:"nutriscore"`
								EtchebestPartnership any    `json:"etchebestPartnership"`
								Type                 string `json:"type"`
								QuitoqueID           any    `json:"quitoqueId"`
								Name                 any    `json:"name"`
								NbMeal               any    `json:"nbMeal"`
								Subtitle             any    `json:"subtitle"`
								Description          any    `json:"description"`
								Components           string `json:"components"`
								StorageWay           string `json:"storageWay"`
								Manual               string `json:"manual"`
								Images               []any  `json:"images"`
								RemainingQuantity    any    `json:"remainingQuantity"`
								QuantityThreshold    any    `json:"quantityThreshold"`
								Responsive           struct {
									ID       string `json:"id"`
									Medium   any    `json:"medium"`
									Large    any    `json:"large"`
									Typename string `json:"__typename"`
								} `json:"responsive"`
								NbPerson               any `json:"nbPerson"`
								PriceIncludingTax      any `json:"priceIncludingTax"`
								OriginalPrice          any `json:"originalPrice"`
								Origin                 any `json:"origin"`
								Weight                 any `json:"weight"`
								Unit                   any `json:"unit"`
								IsFoodProduct          any `json:"isFoodProduct"`
								MaxQuantity            any `json:"maxQuantity"`
								NutritionalInformation struct {
									NbPerson          any    `json:"nbPerson"`
									Protein           any    `json:"protein"`
									Fat               any    `json:"fat"`
									SaturatedFat      any    `json:"saturatedFat"`
									Carbohydrate      any    `json:"carbohydrate"`
									SugarCarbohydrate any    `json:"sugarCarbohydrate"`
									WeightPerPerson   any    `json:"weightPerPerson"`
									Fiber             any    `json:"fiber"`
									Salt              any    `json:"salt"`
									KiloCalorie       any    `json:"kiloCalorie"`
									KiloJoule         any    `json:"kiloJoule"`
									Typename          string `json:"__typename"`
								} `json:"nutritionalInformation"`
								Allergens            []any  `json:"allergens"`
								AllergenTraces       []any  `json:"allergenTraces"`
								Facets               []any  `json:"facets"`
								BoxDescription       any    `json:"boxDescription"`
								BoxByWeekDescription any    `json:"boxByWeekDescription"`
								CookingTime          any    `json:"cookingTime"`
								WaitingTime          any    `json:"waitingTime"`
								UserMark             any    `json:"userMark"`
								Typename             string `json:"__typename"`
								SubProducts          []any  `json:"subProducts"`
							} `json:"product"`
							Typename string `json:"__typename"`
						} `json:"tools"`
						Ingredients []struct {
							ID              string `json:"id"`
							Name            string `json:"name"`
							Position        int    `json:"position"`
							QuantityToGram  any    `json:"quantityToGram"`
							QuantityToOrder any    `json:"quantityToOrder"`
							LiteralQuantity string `json:"literalQuantity"`
							Quantity        int    `json:"quantity"`
							Product         struct {
								ID                   string   `json:"id"`
								Slug                 any      `json:"slug"`
								OriginalID           string   `json:"originalId"`
								WeekDetailID         any      `json:"weekDetailId"`
								Nutriscore           any      `json:"nutriscore"`
								EtchebestPartnership any      `json:"etchebestPartnership"`
								Type                 string   `json:"type"`
								QuitoqueID           any      `json:"quitoqueId"`
								Name                 string   `json:"name"`
								NbMeal               any      `json:"nbMeal"`
								Subtitle             string   `json:"subtitle"`
								Description          string   `json:"description"`
								Components           string   `json:"components"`
								StorageWay           string   `json:"storageWay"`
								Manual               string   `json:"manual"`
								Images               []string `json:"images"`
								RemainingQuantity    any      `json:"remainingQuantity"`
								QuantityThreshold    any      `json:"quantityThreshold"`
								Responsive           struct {
									ID       string `json:"id"`
									Medium   string `json:"medium"`
									Large    string `json:"large"`
									Typename string `json:"__typename"`
								} `json:"responsive"`
								NbPerson               any    `json:"nbPerson"`
								PriceIncludingTax      any    `json:"priceIncludingTax"`
								OriginalPrice          any    `json:"originalPrice"`
								Origin                 string `json:"origin"`
								Weight                 int    `json:"weight"`
								Unit                   string `json:"unit"`
								IsFoodProduct          bool   `json:"isFoodProduct"`
								MaxQuantity            int    `json:"maxQuantity"`
								NutritionalInformation struct {
									NbPerson          any     `json:"nbPerson"`
									Protein           float64 `json:"protein"`
									Fat               float64 `json:"fat"`
									SaturatedFat      float64 `json:"saturatedFat"`
									Carbohydrate      float64 `json:"carbohydrate"`
									SugarCarbohydrate float64 `json:"sugarCarbohydrate"`
									WeightPerPerson   any     `json:"weightPerPerson"`
									Fiber             float64 `json:"fiber"`
									Salt              float64 `json:"salt"`
									KiloCalorie       float64 `json:"kiloCalorie"`
									KiloJoule         float64 `json:"kiloJoule"`
									Typename          string  `json:"__typename"`
								} `json:"nutritionalInformation"`
								Allergens []struct {
									ID       string `json:"id"`
									Name     string `json:"name"`
									Typename string `json:"__typename"`
								} `json:"allergens"`
								AllergenTraces []struct {
									ID       string `json:"id"`
									Name     string `json:"name"`
									Typename string `json:"__typename"`
								} `json:"allergenTraces"`
								Facets               []any  `json:"facets"`
								BoxDescription       any    `json:"boxDescription"`
								BoxByWeekDescription any    `json:"boxByWeekDescription"`
								CookingTime          any    `json:"cookingTime"`
								WaitingTime          any    `json:"waitingTime"`
								UserMark             any    `json:"userMark"`
								Typename             string `json:"__typename"`
								SubProducts          []any  `json:"subProducts"`
							} `json:"product"`
							Typename string `json:"__typename"`
						} `json:"ingredients"`
						CupboardIngredients []struct {
							ID              string `json:"id"`
							Name            string `json:"name"`
							Position        int    `json:"position"`
							QuantityToGram  any    `json:"quantityToGram"`
							QuantityToOrder any    `json:"quantityToOrder"`
							LiteralQuantity string `json:"literalQuantity"`
							Quantity        int    `json:"quantity"`
						} `json:"cupboardIngredients"`
						Typename string `json:"__typename"`
					} `json:"stacks"`
					Steps []struct {
						Position    int    `json:"position"`
						Title       string `json:"title"`
						Description string `json:"description"`
					} `json:"steps"`
				} `json:"cookingModes"`
			} `json:"pools"`
			NutritionalInformations []struct {
				NbPerson          int     `json:"nbPerson"`
				Protein           float64 `json:"protein"`
				Fat               float64 `json:"fat"`
				SaturatedFat      float64 `json:"saturatedFat"`
				Carbohydrate      float64 `json:"carbohydrate"`
				SugarCarbohydrate float64 `json:"sugarCarbohydrate"`
				WeightPerPerson   int     `json:"weightPerPerson"`
				Fiber             float64 `json:"fiber"`
				Salt              float64 `json:"salt"`
				KiloCalorie       float64 `json:"kiloCalorie"`
				KiloJoule         float64 `json:"kiloJoule"`
				Typename          string  `json:"__typename"`
			} `json:"nutritionalInformations"`
		} `json:"recipe"`
	} `json:"data"`
}

func (s *Scraper) scrapeQuitoque(rawURL string) (models.RecipeSchema, error) {
	_, after, ok := strings.Cut(rawURL, "quitoque.fr/recette/")
	if !ok {
		return models.RecipeSchema{}, errors.New("url is invalid")
	}

	before, _, ok := strings.Cut(after, "/")
	if !ok {
		return models.RecipeSchema{}, errors.New("url is invalid")
	}

	api := `https://mgs.quitoque.fr/graphql?operationName=getRecipe&variables={"id":"` + before + `"}&extensions={"persistedQuery":{"version":1,"sha256Hash":"04af4d1a48fd536a67292733e23a2afcf6d0da9770ab07055c59b754eec9bd6d"}}`
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	res, err := s.Client.Do(req)
	if err != nil {
		return models.RecipeSchema{}, err
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var q quitoque
	err = json.Unmarshal(buf, &q)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var ns models.NutritionSchema
	if len(q.Data.Recipe.NutritionalInformations) > 0 {
		n := q.Data.Recipe.NutritionalInformations[0]
		ns = models.NutritionSchema{
			Calories:      strconv.FormatFloat(n.KiloCalorie, 'g', 10, 64) + " kcal",
			Carbohydrates: strconv.FormatFloat(n.Carbohydrate, 'g', 10, 64) + " g",
			Fat:           strconv.FormatFloat(n.Fat, 'g', 10, 64) + " g",
			Fiber:         strconv.FormatFloat(n.Fiber, 'g', 10, 64) + " g",
			Protein:       strconv.FormatFloat(n.Protein, 'g', 10, 64) + " g",
			SaturatedFat:  strconv.FormatFloat(n.SaturatedFat, 'g', 10, 64) + " g",
			Servings:      strconv.Itoa(n.NbPerson),
			Sodium:        strconv.FormatFloat(n.Salt, 'g', 10, 64) + " g",
			Sugar:         strconv.FormatFloat(n.SugarCarbohydrate, 'g', 10, 64) + " g",
		}
	}

	var (
		cook         string
		prep         string
		ingredients  []string
		instructions []string
		yield        int16
	)
	if len(q.Data.Recipe.Pools) > 0 && len(q.Data.Recipe.Pools[0].CookingModes) > 0 {
		m := q.Data.Recipe.Pools[0].CookingModes[0]

		yield = int16(q.Data.Recipe.Pools[0].NbPerson)

		prep = "PT" + strconv.Itoa(m.WaitingTime) + "M"
		cook = "PT" + strconv.Itoa(m.CookingTime) + "M"

		instructions = make([]string, 0, len(m.Steps))
		for _, step := range m.Steps {
			ins := step.Title + "\n" + step.Description
			ins = strings.ReplaceAll(ins, "\u00a0", " ")
			instructions = append(instructions, ins)
		}

		ingredients = make([]string, 0, len(m.Stacks.Ingredients)+len(m.Stacks.Ingredients))
		for _, ing := range m.Stacks.Ingredients {
			var sb strings.Builder
			if ing.LiteralQuantity != "" {
				sb.WriteString(ing.LiteralQuantity)
				sb.WriteString(" ")
			}

			sb.WriteString(ing.Name)
			ingredients = append(ingredients, sb.String())
		}

		for _, ing := range m.Stacks.CupboardIngredients {
			var sb strings.Builder
			if ing.LiteralQuantity != "0" {
				sb.WriteString(ing.LiteralQuantity)
				sb.WriteString(" ")
			}

			sb.WriteString(ing.Name)
			ingredients = append(ingredients, sb.String())
		}
	}

	return models.RecipeSchema{
		AtContext:       atContext,
		AtType:          models.SchemaType{Value: "Recipe"},
		Category:        models.Category{},
		CookTime:        cook,
		CookingMethod:   models.CookingMethod{},
		Cuisine:         models.Cuisine{},
		DateCreated:     "",
		DateModified:    "",
		DatePublished:   "",
		Description:     models.Description{Value: q.Data.Recipe.ShortDescription},
		Keywords:        models.Keywords{},
		Image:           models.Image{Value: q.Data.Recipe.Image},
		Ingredients:     models.Ingredients{Values: ingredients},
		Instructions:    models.Instructions{Values: instructions},
		Name:            q.Data.Recipe.Name,
		NutritionSchema: ns,
		PrepTime:        prep,
		Tools:           models.Tools{},
		Yield:           models.Yield{Value: yield},
	}, nil
}
