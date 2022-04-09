package scraper

import (
	"strconv"
	"strings"
	"time"

	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/regex"
	"golang.org/x/net/html"
)

func scrapeMonsieurCuisine(root *html.Node) (models.RecipeSchema, error) {
	content := getElement(root, "class", "row recipe--header")

	chName := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chName <- s
		}()

		h1 := traverseAll(content, func(node *html.Node) bool {
			return node.Data == "h1"
		})
		s = h1[0].FirstChild.Data
	}()

	chYield := make(chan int16)
	go func() {
		var yield int16
		defer func() {
			_ = recover()
			chYield <- yield
		}()

		node := getElement(content, "class", "recipe-portions")
		node = getElement(node, "class", "info-label")
		parts := strings.Split(strings.TrimSpace(node.FirstChild.Data), " ")
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err == nil {
				yield = int16(i)
			}
		}
	}()

	chPrep := make(chan string)
	chCook := make(chan string)
	go func() {
		var prep string
		var cook string
		defer func() {
			_ = recover()
			chPrep <- prep
			chCook <- cook
		}()

		var prepTime int
		node := getElement(content, "class", "recipe-duration")
		node = getElement(node, "class", "info-label")
		s := strings.ReplaceAll(node.FirstChild.Data, "\n", "")
		parts := strings.Split(strings.TrimSpace(s), " ")
		for _, part := range parts {
			var err error
			prepTime, err = strconv.Atoi(part)
			if err == nil {
				prep = "PT" + part + "M"
				break
			}
		}

		node = getElement(content, "class", "recipe-duration-total")
		node = getElement(node, "class", "info-label")
		s = string(regex.HourMinutes.Find([]byte(strings.TrimSpace(node.FirstChild.Data))))
		parts = strings.Split(s, ":")
		if len(parts) != 2 {
			return
		}
		hoursStr, minutesStr := parts[0], parts[1]
		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return
		}
		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			return
		}
		cookTime := (hours*60 + minutes) - prepTime
		cook = "PT" + strconv.Itoa(cookTime) + "M"
	}()

	chDatePublished := make(chan string)
	chDateCreated := make(chan string)
	chDateModified := make(chan string)
	go func() {
		var pub string
		var created string
		var mod string
		defer func() {
			_ = recover()
			chDatePublished <- pub
			chDateCreated <- created
			chDateModified <- mod
		}()

		node := getElement(content, "class", "recipe-crdate")
		if node != nil {
			s := strings.TrimSpace(node.FirstChild.Data)
			parts := strings.Split(s, " ")
			for _, part := range parts {
				t, err := time.Parse("02.01.2006", part)
				if err == nil {
					f := t.Format("2006-01-02")
					created = f
					pub = f
					break
				}
			}
		}

		node = getElement(content, "class", "recipe-last-edit")
		s := strings.TrimSpace(node.FirstChild.Data)
		parts := strings.Split(s, " ")
		for _, part := range parts {
			t, err := time.Parse("02.01.2006", part)
			if err == nil {
				mod = t.Format("2006-01-02")
			}
		}
	}()

	chImage := make(chan string)
	go func() {
		var s string
		defer func() {
			_ = recover()
			chImage <- s
		}()

		node := getElement(root, "class", "flexed-image-preview")
		s = "https://www.monsieur-cuisine.com" + getAttr(node.FirstChild.FirstChild, "src")
	}()

	chIngredients := make(chan models.Ingredients)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chIngredients <- models.Ingredients{Values: vals}
		}()

		node := getElement(root, "class", "recipe--ingredients-html-item col-md-8")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					vals = append(vals, strings.TrimSpace(c.FirstChild.Data))
				}
				break
			}
		}
	}()

	chInstructions := make(chan models.Instructions)
	go func() {
		var vals []string
		defer func() {
			_ = recover()
			chInstructions <- models.Instructions{Values: vals}
		}()

		node := getElement(root, "class", "recipe--instructions")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ol" {
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					vals = append(vals, strings.TrimSpace(c.FirstChild.Data))
				}
				break
			}
		}
	}()

	chNutrition := make(chan models.NutritionSchema)
	go func() {
		var (
			calories      string
			carbohydrates string
			fat           string
			protein       string
		)
		defer func() {
			_ = recover()
			chNutrition <- models.NutritionSchema{
				Calories:      calories,
				Carbohydrates: carbohydrates,
				Fat:           fat,
				Protein:       protein,
			}
		}()

		node := getElement(root, "class", "recipe--nutrients-html")
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "ul" {
				for c := c.FirstChild; c != nil; c = c.NextSibling {
					parts := strings.Split(c.FirstChild.Data, " ")
					value := strings.Join(parts[1:], " ")

					switch strings.ReplaceAll(parts[0], ":", "") {
					case "calorific":
						calories = c.FirstChild.NextSibling.NextSibling.Data
					case "protein":
						protein = value
					case "fat":
						fat = value
					case "carbohydrates":
						carbohydrates = value
					}
				}
				break
			}
		}
	}()

	return models.RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          models.SchemaType{Value: "Recipe"},
		Name:            <-chName,
		PrepTime:        <-chPrep,
		CookTime:        <-chCook,
		Yield:           models.Yield{Value: <-chYield},
		DatePublished:   <-chDatePublished,
		DateCreated:     <-chDateCreated,
		DateModified:    <-chDateModified,
		Image:           models.Image{Value: <-chImage},
		Ingredients:     <-chIngredients,
		Instructions:    <-chInstructions,
		NutritionSchema: <-chNutrition,
	}, nil
}
