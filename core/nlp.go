package core

import (
	"strings"

	"github.com/jdkato/prose/v2"
	"github.com/reaper47/recipya/data"
)

// NlpExtractIngredients extracts the ingredients from a list of ingredients
// using natural language processing.
//
// For example, if the ingredient is "1 cup of butter", then the result would be "butter".
func NlpExtractIngredients(ingredients []string) []string {
	blacklist := data.Data.ReadBlacklistIngredients()
	tags := map[string]int8{
		"JJ":   0,
		"JJR":  0,
		"JJS":  0,
		"NN":   0,
		"NNP":  0,
		"NNS":  0,
		"NNPS": 0,
	}

	processed := make([]string, len(ingredients))
	for i, ingredient := range ingredients {
		doc, err := prose.NewDocument(ingredient)
		if err != nil {
			continue
		}

		item := ""
		for _, tok := range doc.Tokens() {
			if _, found := tags[tok.Tag]; found {
				s := strings.ToLower(tok.Text)

				if _, found = blacklist[s[:len(s)-1]]; found {
					continue
				}

				if _, found := blacklist[s]; found {
					continue
				}

				if _, found := blacklist[s+"s"]; found {
					continue
				}

				item += s + " "
			}

		}
		processed[i] = strings.TrimSpace(item)
	}
	return processed
}
