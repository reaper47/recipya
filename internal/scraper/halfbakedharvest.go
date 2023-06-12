package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
	"strings"
)

func scrapeHalfBakedHarvest(root *goquery.Document) (models.RecipeSchema, error) {
	rs, err := parseGraph(root)
	if err != nil {
		return models.RecipeSchema{}, err
	}

	var instructions []string
	allInstructions := rs.Instructions.Values[0]
	i := 1
	for {
		start := fmt.Sprintf("%d. ", i)
		end := fmt.Sprintf("%d. ", i+1)

		startIndex := strings.Index(allInstructions, start)
		endIndex := strings.Index(allInstructions[startIndex:], end) + startIndex

		if endIndex < startIndex {
			s := strings.TrimSpace(allInstructions[startIndex+len(start):])
			instructions = append(instructions, s)
			break
		}

		if startIndex >= 0 && endIndex >= startIndex {
			s := strings.TrimSpace(allInstructions[startIndex+len(start) : endIndex])
			instructions = append(instructions, s)
		}
		i++
	}
	rs.Instructions.Values = instructions

	return rs, nil
}
