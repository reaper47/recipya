package models

import (
	"github.com/jdkato/prose/v2"
	"github.com/reaper47/recipya/internal/utils/regex"
	"slices"
	"strconv"
	"strings"
	"time"
)

// AzureDILayout holds the data contained in the response of a call to the Azure AI Document Intelligence layout analysis endpoint.
type AzureDILayout struct {
	Status              string    `json:"status"`
	CreatedDateTime     time.Time `json:"createdDateTime"`
	LastUpdatedDateTime time.Time `json:"lastUpdatedDateTime"`
	AnalyzeResult       struct {
		ApiVersion      string `json:"apiVersion"`
		ModelId         string `json:"modelId"`
		StringIndexType string `json:"stringIndexType"`
		Content         string `json:"content"`
		Pages           []struct {
			PageNumber int     `json:"pageNumber"`
			Angle      float64 `json:"angle"`
			Width      int     `json:"width"`
			Height     int     `json:"height"`
			Unit       string  `json:"unit"`
			Words      []struct {
				Content    string  `json:"content"`
				Polygon    []int   `json:"polygon"`
				Confidence float64 `json:"confidence"`
				Span       struct {
					Offset int `json:"offset"`
					Length int `json:"length"`
				} `json:"span"`
			} `json:"words"`
			Lines []struct {
				Content string `json:"content"`
				Polygon []int  `json:"polygon"`
				Spans   []struct {
					Offset int `json:"offset"`
					Length int `json:"length"`
				} `json:"spans"`
			} `json:"lines"`
			Spans []struct {
				Offset int `json:"offset"`
				Length int `json:"length"`
			} `json:"spans"`
		} `json:"pages"`
		Tables     []interface{}      `json:"tables"`
		Paragraphs []azureDIParagraph `json:"paragraphs"`
		Styles     []struct {
			Confidence float64 `json:"confidence"`
			Spans      []struct {
				Offset int `json:"offset"`
				Length int `json:"length"`
			} `json:"spans"`
			IsHandwritten bool `json:"isHandwritten"`
		} `json:"styles"`
		ContentFormat string `json:"contentFormat"`
		Sections      []struct {
			Spans []struct {
				Offset int `json:"offset"`
				Length int `json:"length"`
			} `json:"spans"`
			Elements []string `json:"elements"`
		} `json:"sections"`
	} `json:"analyzeResult"`
}

type azureDIParagraph struct {
	Spans []struct {
		Offset int `json:"offset"`
		Length int `json:"length"`
	} `json:"spans"`
	BoundingRegions []struct {
		PageNumber int   `json:"pageNumber"`
		Polygon    []int `json:"polygon"`
	} `json:"boundingRegions"`
	Role    string `json:"role,omitempty"`
	Content string `json:"content"`
}

// Recipe converts an AzureDILayout to a Recipe.
func (a *AzureDILayout) Recipe() Recipe {
	recipe := Recipe{
		Category: "uncategorized",
		URL:      "OCR",
		Yield:    1,
	}

	if len(a.AnalyzeResult.Pages) == 0 {
		return Recipe{}
	}

	var (
		isFromBook     = len(a.AnalyzeResult.Styles) == 0 || !a.AnalyzeResult.Styles[0].IsHandwritten
		isMetaData     = true
		isIngredients  bool
		isInstructions bool
	)

	for _, p := range a.AnalyzeResult.Paragraphs {
		if p.Role == "title" {
			recipe.Name = p.Content
			if isFromBook {
				continue
			}
			break
		} else if recipe.Name == "" && p.Role == "sectionHeading" {
			recipe.Name = p.Content
			isFromBook = true
			continue
		} else if p.Role == "pageHeader" || p.Role == "pageNumber" {
			continue
		} else if isInstructions && (p.Role == "sectionHeading" || (p.Role == "" && len(p.Content) > 5)) {
			s := p.Content
			dotIndex := strings.IndexByte(s, '.')
			if dotIndex != -1 && dotIndex < 4 {
				_, s, _ = strings.Cut(s, ".")
			}
			recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(s))
			continue
		} else if len(p.Content) < 3 {
			continue
		}

		if isFromBook {
			if isMetaData {
				if len(p.Content) < 20 && strings.Contains(strings.ToLower(p.Content), "serve") {
					parsed, err := strconv.ParseInt(regex.Digit.FindString(p.Content), 10, 16)
					if err == nil {
						recipe.Yield = int16(parsed)
					}
				} else if !isIngredient(p.Content) {
					recipe.Description = p.Content
					isMetaData = false
					isIngredients = true
				}
			} else if isIngredients {
				isIng := isIngredient(p.Content)

				if !isIng {
					s := p.Content
					dotIndex := strings.IndexByte(s, '.')
					if dotIndex < 4 {
						_, s, _ = strings.Cut(s, ".")
					}
					recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(s))

					isIngredients = false
					isInstructions = true
				} else {
					recipe.Ingredients = append(recipe.Ingredients, p.Content)
				}
			}
		}
	}

	if isFromBook {
		return recipe
	}

	isInstructionBlock := true
	for _, page := range a.AnalyzeResult.Pages {
		for _, line := range page.Lines {
			if line.Content == recipe.Name {
				continue
			} else if !isIngredient(line.Content) {
				goto Break
			} else if len(recipe.Ingredients) == 0 {
				doc, _ := prose.NewDocument(line.Content)
				tokens := doc.Tokens()
				if len(tokens) > 0 && tokens[0].Tag == "IN" {
					recipe.Description = line.Content
					continue
				}
			}

			recipe.Ingredients = append(recipe.Ingredients, line.Content)
		}
	}

Break:
	pidx := slices.IndexFunc(a.AnalyzeResult.Paragraphs, func(xp azureDIParagraph) bool {
		return strings.Contains(xp.Content, recipe.Ingredients[len(recipe.Ingredients)-1])
	})
	if pidx != -1 {
		for _, p := range a.AnalyzeResult.Paragraphs[pidx+1:] {
			if p.Role != "" {
				continue
			}

			if len(p.Content) < 5 {
				continue
			}

			_, after, ok := strings.Cut(strings.ToLower(p.Content), "serve")
			if ok && regex.Digit.MatchString(after) {
				isInstructionBlock = false
				parsed, err := strconv.ParseInt(regex.Digit.FindString(p.Content), 10, 16)
				if err == nil {
					recipe.Yield = int16(parsed)
					continue
				}
			}

			if isInstructionBlock {
				recipe.Instructions = append(recipe.Instructions, p.Content)
				continue
			}
		}
	}

	recipe.Description = strings.TrimSpace(recipe.Description)
	if recipe.Description == "" {
		recipe.Description = "Recipe created using Azure AI Document Intelligence."

	}
	return recipe
}

// AzureDIError holds the data of an Azure AI Document Intelligence error.
type AzureDIError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func isIngredient(s string) bool {
	if s == "" {
		return false
	}

	_, err := strconv.ParseFloat(string([]rune(s)[0]), 64)
	dotIndex := strings.Index(s, ".")

	idx := regex.Unit.FindStringIndex(s)
	isGood := idx != nil && idx[0] < 8

	return isGood || (err == nil && (dotIndex == -1 || dotIndex > 3)) || len(s) < 25
}
