package models

import (
	"github.com/jdkato/prose/v2"
	"slices"
	"strconv"
	"strings"
)

// AzureVision holds the data contained in the response of a call to the
// Computer Vision API's image analysis endpoint.
type AzureVision struct {
	CaptionResult struct {
		Text       string  `json:"text"`
		Confidence float64 `json:"confidence"`
	} `json:"captionResult"`
	ReadResult struct {
		StringIndexType string `json:"stringIndexType"`
		Content         string `json:"content"`
		Pages           []struct {
			Height     float64 `json:"height"`
			Width      float64 `json:"width"`
			Angle      float64 `json:"angle"`
			PageNumber int     `json:"pageNumber"`
			Words      []struct {
				Content     string    `json:"content"`
				BoundingBox []float64 `json:"boundingBox"`
				Confidence  float64   `json:"confidence"`
				Span        struct {
					Offset int `json:"offset"`
					Length int `json:"length"`
				} `json:"span"`
			} `json:"words"`
			Spans []struct {
				Offset int `json:"offset"`
				Length int `json:"length"`
			} `json:"spans"`
			Lines []struct {
				Content     string    `json:"content"`
				BoundingBox []float64 `json:"boundingBox"`
				Spans       []struct {
					Offset int `json:"offset"`
					Length int `json:"length"`
				} `json:"spans"`
			} `json:"lines"`
		} `json:"pages"`
		Styles []struct {
			IsHandwritten bool `json:"isHandwritten"`
			Spans         []struct {
				Offset int `json:"offset"`
				Length int `json:"length"`
			} `json:"spans"`
			Confidence float64 `json:"confidence"`
		} `json:"styles"`
		ModelVersion string `json:"modelVersion"`
	} `json:"readResult"`
	ModelVersion string `json:"modelVersion"`
	Metadata     struct {
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"metadata"`
}

// Recipe converts an AzureVision to a Recipe.
func (a *AzureVision) Recipe() Recipe {
	r := Recipe{
		Category:    "uncategorized",
		Description: a.CaptionResult.Text,
		URL:         a.CaptionResult.Text,
		Yield:       4,
	}

	xs := strings.Split(a.ReadResult.Content, "\n")
	if len(xs) == 0 {
		return r
	}

	r.Name = xs[0]
	var instructionsStartIndex int

loop:
	for i, s := range xs[1:] {
		doc, err := prose.NewDocument(s)
		if err != nil {
			continue
		}

		tokens := doc.Tokens()
		isAllOutside := true
		var hasCD bool

		_, err = strconv.ParseFloat(tokens[0].Text, 64)
		if err == nil {
			hasCD = true
		}

		for _, t := range tokens {
			if t.Tag == "CD" {
				hasCD = true
			}

			if !hasCD && ((len(tokens) > 4 && !hasCD && i > 1) || (t.Label != "O" && i > 2 && len(tokens) > 2)) {
				isAllOutside = false
				instructionsStartIndex = i + 1
				if hasCD {
					instructionsStartIndex++
					r.Ingredients = append(r.Ingredients, doc.Text)
					break
				}
				break loop
			}
		}

		if isAllOutside || hasCD {
			r.Ingredients = append(r.Ingredients, doc.Text)
		}
	}

	if instructionsStartIndex < len(xs) {
		joined := xs[instructionsStartIndex:]
		var (
			instructions *prose.Document
			err          error
		)
		if slices.ContainsFunc(joined, func(s string) bool { return strings.Contains(s, ".") }) {
			instructions, err = prose.NewDocument(strings.Join(joined, " "))
		} else {
			instructions, err = prose.NewDocument(strings.Join(joined, ". "))
		}

		if err != nil {
			return r
		}

		sentences := instructions.Sentences()
		r.Instructions = make([]string, len(sentences))
		for i, s := range sentences {
			r.Instructions[i] = s.Text
		}
	}

	return r
}
