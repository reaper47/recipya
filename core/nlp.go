package core

import (
	"log"
	"strings"

	"github.com/jdkato/prose/v2"
)

// NlpExtractIngredients extracts the ingredients from a list of ingredients
// using natural language processing.
//
// For example, if the ingredient is "1 cup of butter", then the result would be "butter".
func (env *Env) NlpExtractIngredients(ingredients []string) []string {
	blacklist, err := env.data.GetBlacklistUnits()
	if err != nil {
		log.Printf("Error getting blacklisted ingredients: '%s'\n", err)
		return []string{}
	}

	produce, err := env.data.GetFruitsVeggies()
	if err != nil {
		log.Printf("Error getting fruits and vegetables: '%s'\n", err)
		return []string{}
	}

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

		item, isPlural := parseTokens(doc, tags, blacklist)
		item = strings.TrimSpace(item)
		if isNaturalProduce(item, isPlural, produce) {
			item += " raw"
		}
		processed[i] = item
	}
	return processed
}

func parseTokens(doc *prose.Document, tags, blacklist map[string]int8) (string, bool) {
	item := ""
	var isPlural bool
	for _, tok := range doc.Tokens() {
		if tok.Tag == "NNS" || tok.Tag == "NNPS" {
			isPlural = true
		}

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

			if _, found := blacklist[s+"s"]; found {
				continue
			}
			item += s + " "
		}
	}
	return item, isPlural
}

func isNaturalProduce(item string, isPlural bool, produce map[string]int8) bool {
	for _, word := range strings.Split(item, " ") {
		if isPlural {
			word = getSingular(word)
		}

		if _, found := produce[word]; found {
			return true
		}
	}
	return false
}

func getSingular(word string) string {
	if word[len(word)-3:] == "ies" {
		return word[:len(word)-3] + "y"
	}

	if word[len(word)-2:] == "es" && word[len(word)-3] == 'h' {
		return word[:len(word)-2]
	}

	if word == "tomatoes" || word == "potatoes" {
		return word[:len(word)-2]
	}

	if word[len(word)-1] == 's' {
		return word[:len(word)-1]
	}

	return word
}
