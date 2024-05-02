package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"mime/multipart"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type cookML struct {
	XMLName xml.Name `xml:"cookml"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Prog    string   `xml:"prog,attr"`
	ProgVer string   `xml:"progver,attr"`
	Recipe  []struct {
		Text string `xml:",chardata"`
		Lang string `xml:"lang,attr"`
		Head struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title,attr"`
			Rid         string `xml:"rid,attr"`
			ServingQty  string `xml:"servingqty,attr"`
			ServingType string `xml:"servingtype,attr"`
			Quality     string `xml:"quality,attr"`
			Difficulty  string `xml:"difficulty,attr"`
			CreateDate  string `xml:"createdate,attr"`
			CreateUser  string `xml:"createuser,attr"`
			ChangeDate  string `xml:"changedate,attr"`
			ChangeUser  string `xml:"changeuser,attr"`
			TimeAllQty  string `xml:"timeallqty,attr"`
			CreateEmail string `xml:"createemail,attr"`
			Costs       string `xml:"costs,attr"`
			Cat         []struct {
				Text string `xml:",chardata"`
			} `xml:"cat"`
			Hint []struct {
				Text string `xml:",chardata"`
			} `xml:"hint"`
			SourceLine []struct {
				Text string `xml:",chardata"`
			} `xml:"sourceline"`
			PicBin struct {
				Text   string `xml:",chardata"`
				Format string `xml:"format,attr"`
			} `xml:"picbin"`
		} `xml:"head"`
		Part []struct {
			Text       string `xml:",chardata"`
			Title      string `xml:"title,attr"`
			Ingredient []struct {
				Text string `xml:",chardata"`
				Qty  string `xml:"qty,attr"`
				Unit string `xml:"unit,attr"`
				Item string `xml:"item,attr"`
				Bls  string `xml:"bls,attr"`
				Gram string `xml:"gram,attr"`
				Shop string `xml:"shop,attr"`
				Note []struct {
					Text string `xml:",chardata"`
				} `xml:"inote"`
			} `xml:"ingredient"`
		} `xml:"part"`
		Preparation struct {
			Chardata string `xml:",chardata"`
			Text     struct {
				Text string `xml:",chardata"`
			} `xml:"text"`
		} `xml:"preparation"`
		Remark struct {
			Text string `xml:",chardata"`
			User string `xml:"user,attr"`
			Line struct {
				Text string `xml:",chardata"`
			} `xml:"line"`
		} `xml:"remark"`
	} `xml:"recipe"`
}

type crouton struct {
	Name            string        `json:"name"`
	WebLink         string        `json:"webLink"`
	UUID            string        `json:"uuid"`
	CookingDuration int           `json:"cookingDuration"`
	Images          []string      `json:"images"`
	Serves          int           `json:"serves"`
	FolderIDs       []interface{} `json:"folderIDs"`
	DefaultScale    int           `json:"defaultScale"`
	Ingredients     []struct {
		Order      int `json:"order"`
		Ingredient struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"ingredient"`
		UUID string `json:"uuid"`
	} `json:"ingredients"`
	Steps []struct {
		Step      string `json:"step"`
		Order     int    `json:"order"`
		IsSection bool   `json:"isSection"`
		UUID      string `json:"uuid"`
	} `json:"steps"`
	Tags []struct {
		UUID  string `json:"uuid"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"tags"`
	IsPublicRecipe bool   `json:"isPublicRecipe"`
	SourceImage    string `json:"sourceImage"`
	Nutrition      string `json:"neutritionalInfo"`
	SourceName     string `json:"sourceName"`
	Duration       int    `json:"duration"`
}

// NewRecipesFromCML extracts all recipes from a CookML file.
func NewRecipesFromCML(r io.Reader, file *multipart.FileHeader, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error)) Recipes {
	if r == nil {
		openedFile, err := file.Open()
		if err != nil {
			slog.Error("Failed to open file", "file", file, "error", err)
			return nil
		}
		defer openedFile.Close()
		r = openedFile
	}

	data, err := io.ReadAll(r)
	if err != nil {
		slog.Error("Failed to read file", "file", file.Filename, "error", err)
		return nil
	}

	cleanedData := strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return '�'
		}
		return r
	}, string(data))

	var c cookML
	err = xml.NewDecoder(strings.NewReader(cleanedData)).Decode(&c)
	if err != nil {
		slog.Error("Failed to decode CML", "file", file.Filename, "error", err)
		return nil
	}

	recipes := make(Recipes, 0, len(c.Recipe))
	for _, recipe := range c.Recipe {
		var yield int16 = 1
		parsed, err := strconv.ParseInt(recipe.Head.ServingQty, 10, 16)
		if err == nil {
			yield = int16(parsed)
		}

		var dateCreated time.Time
		parts := strings.Split(recipe.Head.CreateDate, "-")
		if len(parts[0]) == 1 {
			parts[0] = "200" + parts[0]
		}
		before, _, ok := strings.Cut(strings.Join(parts, "-"), "T")
		if ok {
			dateCreated, _ = time.Parse(time.DateOnly, before)
		}

		before, _, _ = strings.Cut(recipe.Head.ChangeDate, "T")
		dateMod, _ := time.Parse(time.DateOnly, before)

		keywords := make([]string, 0, len(recipe.Head.Hint))
		for _, h := range recipe.Head.Hint {
			keywords = append(keywords, h.Text)
		}

		var category string
		if len(keywords) > 0 {
			category = keywords[0]
		}

		if category == "" {
			category = "uncategorized"
		}

		var images []uuid.UUID
		decode, err := base64.StdEncoding.DecodeString(recipe.Head.PicBin.Text)
		if err == nil {
			img, err := uploadImageFunc(io.NopCloser(bytes.NewReader(decode)))
			if err != nil {
				slog.Error("Failed to upload CML image", "src", recipe.Head.PicBin.Text, "error", err)
			}

			if img != uuid.Nil {
				images = append(images, img)
			}
		}

		var prep time.Duration
		if recipe.Head.TimeAllQty != "" {
			parts := strings.Split(recipe.Head.TimeAllQty, " ")
			for _, part := range parts {
				d, err := time.ParseDuration(part + "m")
				if err != nil {
					continue
				}
				prep = d
			}
		}

		var ingredients []string
		for _, p := range recipe.Part {
			ingredients = append(ingredients, p.Title)
			for _, ing := range p.Ingredient {
				var s string
				if ing.Qty != "" {
					s += ing.Qty + " "
				}
				if ing.Unit != "" {
					s += ing.Unit + " "
				}
				if ing.Item != "" {
					s += ing.Item
				}

				ingredients = append(ingredients, strings.TrimSpace(s))
			}
		}

		parts = strings.Split(recipe.Preparation.Text.Text, "\n\n")
		instructions := make([]string, 0, len(parts))
		for _, part := range parts {
			dotIndex := strings.IndexByte(part, '.')
			if dotIndex < 4 {
				_, part, _ = strings.Cut(part, ".")
			}
			instructions = append(instructions, strings.ReplaceAll(strings.TrimSpace(part), "\n", " "))
		}

		recipes = append(recipes, Recipe{
			Category:    category,
			CreatedAt:   dateCreated,
			Description: "Imported from a CookML file.",
			Images:      images,
			Ingredients: slices.DeleteFunc(ingredients, func(s string) bool {
				return s == ""
			}),
			Instructions: instructions,
			Keywords:     keywords,
			Name:         recipe.Head.Title,
			Times:        Times{Prep: prep},
			UpdatedAt:    dateMod,
			URL:          "CookML file",
			Yield:        yield,
		})
	}
	return recipes
}

// NewRecipeFromCrouton create a Recipe from the content of a Crouton file.
func NewRecipeFromCrouton(r io.Reader, uploadImageFunc func(rc io.ReadCloser) (uuid.UUID, error)) Recipe {
	var c crouton
	err := json.NewDecoder(r).Decode(&c)
	if err != nil {
		slog.Error("Failed to decode crouton recipe", "error", err)
		return Recipe{}
	}

	src := c.WebLink
	if src != "" {
		src = "Crouton"
	}

	var images []uuid.UUID
	if len(c.Images) > 0 {
		decode, err := base64.StdEncoding.DecodeString(c.Images[0])
		if err == nil {
			img, err := uploadImageFunc(io.NopCloser(bytes.NewReader(decode)))
			if err != nil {
				slog.Error("Failed to upload Crouton image", "src", c.SourceImage, "file", c.Name, "error", err)
			}

			if img != uuid.Nil {
				images = append(images, img)
			}
		}
	}

	ingredients := make([]string, 0, len(c.Ingredients))
	for _, ing := range c.Ingredients {
		ingredients = append(ingredients, ing.Ingredient.Name)
	}

	instructions := make([]string, 0, len(c.Steps))
	for _, step := range c.Steps {
		instructions = append(instructions, step.Step)
	}

	keywords := make([]string, 0, len(c.Tags))
	for _, tag := range c.Tags {
		keywords = append(keywords, tag.Name)
	}

	category := "uncategorized"
	if len(keywords) > 0 {
		category = keywords[0]
	}

	var n Nutrition
	for _, s := range strings.Split(c.Nutrition, ",\n") {
		before, after, ok := strings.Cut(s, ":")
		if !ok {
			continue
		}

		after = strings.TrimSpace(after)

		switch strings.ToLower(before) {
		case "carbohydrates":
			n.TotalCarbohydrates = after
		case "calories":
			n.Calories = after
		case "fat":
			n.TotalFat = after
		case "sugar":
			n.Sugars = after
		case "cholesterol":
			n.Cholesterol = after
		case "fiber":
			n.Fiber = after
		case "saturated fat":
			n.SaturatedFat = after
		case "sodium":
			n.Sodium = after
		case "protein":
			n.Protein = after
		}
	}

	return Recipe{
		Category:     category,
		Description:  "Imported from Crouton",
		Images:       images,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     keywords,
		Name:         c.Name,
		Nutrition:    n,
		Times: Times{
			Prep: time.Duration(c.Duration) * time.Minute,
			Cook: time.Duration(c.CookingDuration) * time.Minute,
		},
		Tools: nil,
		URL:   src,
		Yield: int16(c.Serves),
	}
}
