package models

import (
	"github.com/reaper47/recipya/internal/units"
	"github.com/reaper47/recipya/internal/utils/regex"
	"strconv"
	"strings"
)

// AzureVision holds the data contained in the response of a call to the Computer Vision API's image analysis endpoint.
type AzureVision struct {
	ModelVersion  string `json:"modelVersion"`
	CaptionResult struct {
		Text       string  `json:"text"`
		Confidence float64 `json:"confidence"`
	} `json:"captionResult"`
	Metadata struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"metadata"`
	ReadResult struct {
		Blocks []struct {
			Lines []struct {
				Text            string `json:"text"`
				BoundingPolygon []struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"boundingPolygon"`
				Words []struct {
					Text            string `json:"text"`
					BoundingPolygon []struct {
						X int `json:"x"`
						Y int `json:"y"`
					} `json:"boundingPolygon"`
					Confidence float64 `json:"confidence"`
				} `json:"words"`
			} `json:"lines"`
		} `json:"blocks"`
	} `json:"readResult"`
}

// Recipe converts an AzureVision to a Recipe.
func (a *AzureVision) Recipe() Recipe {
	recipe := Recipe{
		Category: "uncategorized",
		URL:      "OCR",
	}

	if len(a.ReadResult.Blocks) == 0 || (len(a.ReadResult.Blocks) > 0 && len(a.ReadResult.Blocks[0].Lines) <= 2) {
		return Recipe{}
	}

	var (
		isDescription  = true
		isIngredients  bool
		isInstructions bool

		isDelimitedByAsterisk bool
		isDelimitedByNothing  bool
		isDelimitedByNumbers  bool
	)

	for _, line := range a.ReadResult.Blocks[0].Lines {
		if recipe.Name == "" {
			recipe.Name = line.Text
			continue
		}

		lower := strings.ToLower(line.Text)
		numWords := len(strings.Split(lower, " "))
		if numWords < 3 && strings.HasPrefix(lower, "serve") {
			processMetaData([]byte(lower), &recipe)
			continue
		}

		if isDescription || (isDescription && recipe.Description == "") {
			replaced := units.ReplaceVulgarFractions(line.Text)
			idx := regex.Digit.FindStringIndex(replaced)
			if len(idx) == 2 {
				_, err := strconv.ParseInt(string(replaced[idx[1]-1]), 10, 64)
				if err == nil && idx[0] == 0 {
					isDescription = false
					isIngredients = true
					recipe.Ingredients = append(recipe.Ingredients, line.Text)
					continue
				}
			}

			isDescription = true
			recipe.Description += " " + line.Text
			continue
		}

		if isIngredients {
			replaced := units.ReplaceVulgarFractions(line.Text)
			idx := regex.Digit.FindStringIndex(replaced)
			if len(idx) == 2 {
				_, err := strconv.ParseInt(string(replaced[idx[1]-1]), 10, 64)
				if err == nil && idx[0] == 0 {
					recipe.Ingredients = append(recipe.Ingredients, line.Text)
					continue
				}
			}

			dotIndex := strings.IndexByte(line.Text, '.')
			before, _, found := strings.Cut(line.Text, ".")
			if found && dotIndex >= 0 && dotIndex < 4 && len(line.Text) > 25 {
				_, err := strconv.ParseInt(before, 10, 64)
				if err == nil && dotIndex != -1 && dotIndex < 4 {
					isDelimitedByNumbers = true
					isDelimitedByNothing = false
					isIngredients = false
					isInstructions = true
					recipe.Instructions = append(recipe.Instructions, line.Text)
					continue
				}
			}

			if len(line.Text) >= 25 {
				isIngredients = false
				isInstructions = true

				_, err := strconv.ParseInt(line.Text[:1], 10, 64)
				if err != nil {
					if line.Text[0] == '*' {
						isDelimitedByAsterisk = true
					} else {
						isDelimitedByNothing = true
					}
				}

				recipe.Instructions = append(recipe.Instructions, line.Text)
				continue
			}

			recipe.Ingredients = append(recipe.Ingredients, line.Text)
			continue
		}

		if isInstructions {
			replaced := units.ReplaceVulgarFractions(line.Text)
			idx := regex.Digit.FindStringIndex(replaced)
			if len(idx) == 2 {
				_, err := strconv.ParseInt(string(replaced[idx[1]-1]), 10, 64)
				if err == nil && idx[0] == 0 && !regex.Time.MatchString(line.Text) {
					isInstructions = false
					isIngredients = true
					recipe.Ingredients = append(recipe.Ingredients, line.Text)
					continue
				}
			}

			if line.Text[0] == '*' {
				isDelimitedByAsterisk = true
				recipe.Instructions = append(recipe.Instructions, line.Text)
				continue
			}

			if isDelimitedByAsterisk {
				recipe.Instructions[len(recipe.Instructions)-1] += " " + line.Text
				continue
			}

			if isDelimitedByNumbers {

			}

			if isDelimitedByNothing {
				recipe.Instructions = append(recipe.Instructions, line.Text)
				continue
			}

			dotIndex := strings.IndexByte(line.Text, '.')
			before, _, found := strings.Cut(line.Text, ".")
			if found && dotIndex >= 0 && dotIndex < 4 {
				_, err := strconv.ParseInt(before, 10, 64)
				if err == nil && dotIndex != -1 && dotIndex < 4 {
					recipe.Instructions = append(recipe.Instructions, line.Text)
					continue
				}
			}

			recipe.Instructions[len(recipe.Instructions)-1] += " " + line.Text
			continue
		}
	}

	recipe.Description = strings.TrimSpace(recipe.Description)
	if recipe.Yield == 0 {
		recipe.Yield = 1
	}
	return recipe
}
