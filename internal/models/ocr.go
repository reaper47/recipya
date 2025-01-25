package models

import (
	"github.com/jdkato/prose/v2"
	"github.com/reaper47/recipya/internal/utils/duration"
	"github.com/reaper47/recipya/internal/utils/extensions"
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
		APIVersion      string `json:"apiVersion"`
		ModelID         string `json:"modelId"`
		StringIndexType string `json:"stringIndexType"`
		Content         string `json:"content"`
		Pages           []struct {
			PageNumber int     `json:"pageNumber"`
			Angle      float64 `json:"angle"`
			Width      float64 `json:"width"`
			Height     float64 `json:"height"`
			Unit       string  `json:"unit"`
			Words      []struct {
				Content    string    `json:"content"`
				Polygon    []float64 `json:"polygon"`
				Confidence float64   `json:"confidence"`
				Span       struct {
					Offset int `json:"offset"`
					Length int `json:"length"`
				} `json:"span"`
			} `json:"words"`
			Lines []struct {
				Content string    `json:"content"`
				Polygon []float64 `json:"polygon"`
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
		PageNumber int       `json:"pageNumber"`
		Polygon    []float64 `json:"polygon"`
	} `json:"boundingRegions"`
	Role    string `json:"role,omitempty"`
	Content string `json:"content"`
}

// Recipe converts an AzureDILayout to a Recipe.
func (a *AzureDILayout) Recipe() Recipe {
	recipe := NewBaseRecipe()
	recipe.URL = "OCR"

	if len(a.AnalyzeResult.Pages) == 0 {
		return Recipe{}
	}

	for i := 0; i < len(a.AnalyzeResult.Paragraphs); i++ {
		p := a.AnalyzeResult.Paragraphs[i]

		if len(p.Content) == 1 {
			continue
		}

		if p.Role == "sectionHeading" && strings.HasPrefix(strings.ToLower(p.Content), "utens") {
			if i > 1 {
				recipe.Name = a.AnalyzeResult.Paragraphs[i-1].Content

				_, after, ok := strings.Cut(strings.ToUpper(recipe.Name), "ZUBEREITUNG")
				if ok {
					recipe.Times.Prep = duration.From(strings.ToLower(after))
				}

				_, after, ok = strings.Cut(strings.ToUpper(after), "ETWA")
				if ok {
					parsed, err := strconv.ParseInt(regex.Digit.FindString(after), 10, 16)
					if err == nil {
						recipe.Yield = int16(parsed)
					}
				}
			}

			for i2, p2 := range a.AnalyzeResult.Paragraphs[i+1:] {
				if strings.HasPrefix(strings.ToLower(p2.Content), "zutaten") {
					i += i2 + 1
					break
				}
				recipe.Tools = append(recipe.Tools, NewHowToTool(p2.Content))
			}
			continue
		}

		if recipe.Name == "" {
			if p.Role == "title" || p.Role == "sectionHeading" {
				recipe.Name = p.Content
				continue
			} else if p.Role == "pageHeader" {
				header := regex.Digit.ReplaceAllString(p.Content, "")
				header = strings.ReplaceAll(header, ".", "")
				if header != "" && len(p.Content) < 20 {
					recipe.Category = strings.ToLower(strings.TrimSpace(header))
				}
			}
			continue
		}

		if strings.Contains(strings.ToLower(p.Content), "serve") && len(p.Content) < 20 {
			_, after, _ := strings.Cut(strings.ToLower(p.Content), "serve")
			if regex.Digit.MatchString(after) {
				parsed, _ := strconv.ParseInt(regex.Digit.FindString(after), 10, 16)
				recipe.Yield = int16(parsed)
				continue
			}
		}

		if recipe.Cuisine == "" && recipe.Description == "" && len(p.Content) < 20 {
			doc, _ := prose.NewDocument(p.Content)
			if slices.ContainsFunc(doc.Entities(), func(e prose.Entity) bool { return e.Label == "GPE" }) {
				recipe.Cuisine = strings.ToLower(p.Content)
				continue
			}
		}

		if len(recipe.Ingredients) == 0 {
			if isIngredient(p.Content) {
				var isFound bool
				var isSkipped bool
				for _, line := range a.AnalyzeResult.Pages[p.BoundingRegions[0].PageNumber-1].Lines {
					if !isSkipped {
						if len(line.Content) < 2 {
							continue
						} else if strings.Contains(p.Content, line.Content) {
							isSkipped = true
						} else {
							continue
						}
					}

					if len(line.Content) < 2 {
						continue
					} else if strings.Contains(p.Content, line.Content) || strings.HasPrefix(strings.ToLower(line.Content), "for") {
						isFound = true
						recipe.Ingredients = append(recipe.Ingredients, line.Content)
						continue
					} else if !isFound {
						continue
					} else if !isIngredient(line.Content) {
						i = slices.IndexFunc(a.AnalyzeResult.Paragraphs, func(paragraph azureDIParagraph) bool {
							return strings.Contains(paragraph.Content, line.Content)
						}) - 1
						break
					}

					recipe.Ingredients = append(recipe.Ingredients, line.Content)
				}
				continue
			} else if strings.HasPrefix(strings.ToUpper(p.Content), "FÜR DIE") {
				for i2, p2 := range a.AnalyzeResult.Paragraphs[i:] {
					diff := len(p2.Content) - len(recipe.Name)
					if p2.Role == "title" && (diff < -3 || diff > 3) {
						i += i2
						break
					}
					recipe.Ingredients = append(recipe.Ingredients, p2.Content)
				}
				continue
			} else if recipe.Description != "" {
				if strings.HasSuffix(strings.ToLower(p.Content), "servings") {
					parsed, err := strconv.ParseInt(regex.Digit.FindString(p.Content), 10, 16)
					if err == nil {
						recipe.Yield = int16(parsed)
					}

					continue
				}

				recipe.Description += "\n\n" + p.Content
				continue
			}
		}

		if len(recipe.Ingredients) == 0 {
			recipe.Description = p.Content
			continue
		}

		if len(recipe.Instructions) == 0 {
			if isIngredient(p.Content) {
				var isFound bool
				for _, line := range a.AnalyzeResult.Pages[p.BoundingRegions[0].PageNumber-1].Lines {
					if strings.Contains(p.Content, line.Content) {
						isFound = true
						recipe.Ingredients = append(recipe.Ingredients, line.Content)
						continue
					} else if !isFound {
						continue
					} else if !isIngredient(line.Content) {
						i = slices.IndexFunc(a.AnalyzeResult.Paragraphs, func(paragraph azureDIParagraph) bool {
							return strings.Contains(paragraph.Content, line.Content)
						}) - 1
						break
					}

					recipe.Ingredients = append(recipe.Ingredients, line.Content)
				}
				continue
			}

			for i2, p2 := range a.AnalyzeResult.Paragraphs[i:] {
				if len(strings.Split(p2.Content, " ")) < 2 {
					_, err := strconv.ParseInt(p2.Content, 10, 64)
					if err == nil {
						continue
					}
					i += i2
					break
				} else if strings.Contains(strings.ToLower(p2.Content), "serve") && len(p2.Content) < 20 {
					_, after, _ := strings.Cut(strings.ToLower(p2.Content), "serve")
					if regex.Digit.MatchString(after) {
						parsed, _ := strconv.ParseInt(regex.Digit.FindString(after), 10, 16)
						recipe.Yield = int16(parsed)
						i += i2
						break
					}
				}

				s := p2.Content
				dotIndex := strings.IndexByte(s, '.')
				if dotIndex != -1 && dotIndex < 4 {
					_, s, _ = strings.Cut(s, ".")
				}

				recipe.Instructions = append(recipe.Instructions, strings.TrimSpace(s))
			}
			continue
		}

		if strings.Contains(strings.ToLower(p.Content), "serve") {
			parsed, err := strconv.ParseInt(regex.Digit.FindString(p.Content), 10, 16)
			if err == nil {
				recipe.Yield = int16(parsed)
			}
		}
	}

	if recipe.Description == "" {
		recipe.Description = "Recipe created using Azure AI Document Intelligence."
	}

	// Transfer ingredients in the instructions to ingredients.
	for _, ing := range recipe.Instructions {
		if isIngredient(ing) || strings.HasPrefix(strings.ToLower(ing), "for") {
			recipe.Ingredients = append(recipe.Ingredients, ing)
		} else {
			break
		}
	}

	elements := make(map[string]bool)
	for _, v := range recipe.Ingredients {
		elements[v] = true
	}

	result := make([]string, 0)
	for _, v := range recipe.Instructions {
		if !elements[v] {
			result = append(result, v)
		}
	}

	if len(result) == 0 {
		result = append(result, "No instructions found in image.")
	}

	recipe.Ingredients = extensions.Unique(recipe.Ingredients)
	recipe.Instructions = result
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
	if s == "" || strings.Contains(strings.ToLower(s), "serving") {
		return false
	}

	_, err := strconv.ParseFloat(string([]rune(s)[0]), 64)
	dotIndex := strings.Index(s, ".")

	idx := regex.Unit.FindStringIndex(s)
	isAtStart := idx != nil && idx[0] < 8

	if idx == nil {
		doc, _ := prose.NewDocument(s)
		if doc.Tokens()[0].Tag == "IN" ||
			slices.ContainsFunc(doc.Tokens(), func(e prose.Token) bool { return strings.Contains(e.Tag, "GPE") }) {
			return false
		}
	}

	return isAtStart || (err == nil && (dotIndex == -1 || dotIndex > 3)) || len(s) < 25
}
