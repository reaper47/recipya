package models

import (
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/utils/extensions"
	"log/slog"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecipeSchema is a representation of the Recipe schema (https://schema.org/Recipe).
type RecipeSchema struct {
	AtContext       string           `json:"@context"`
	AtGraph         []*RecipeSchema  `json:"@graph,omitempty"`
	AtType          *SchemaType      `json:"@type"`
	Category        *Category        `json:"recipeCategory,omitempty"`
	CookTime        string           `json:"cookTime,omitempty"`
	CookingMethod   *CookingMethod   `json:"cookingMethod,omitempty"`
	Cuisine         *Cuisine         `json:"recipeCuisine,omitempty"`
	DateCreated     string           `json:"dateCreated,omitempty"`
	DateModified    string           `json:"dateModified,omitempty"`
	DatePublished   string           `json:"datePublished,omitempty"`
	Description     *Description     `json:"description"`
	Keywords        *Keywords        `json:"keywords,omitempty"`
	Image           *Image           `json:"image,omitempty"`
	Ingredients     *Ingredients     `json:"recipeIngredient,omitempty"`
	Instructions    *Instructions    `json:"recipeInstructions,omitempty"`
	Name            string           `json:"name,omitempty"`
	NutritionSchema *NutritionSchema `json:"nutrition,omitempty"`
	PrepTime        string           `json:"prepTime,omitempty"`
	ThumbnailURL    *ThumbnailURL    `json:"thumbnailUrl,omitempty"`
	Tools           *Tools           `json:"tool,omitempty"`
	TotalTime       string           `json:"totalTime,omitempty"`
	Yield           *Yield           `json:"recipeYield,omitempty"`
	URL             string           `json:"url,omitempty"`
	Video           *Videos          `json:"video,omitempty"`
}

// Equal verifies whether a RecipeSchema is equal to the other.
func (r *RecipeSchema) Equal(other RecipeSchema) bool {
	return r.AtType != nil && r.AtType.Value == other.AtType.Value &&
		r.Category != nil && r.Category.Value == other.Category.Value &&
		r.CookTime == other.CookTime &&
		r.CookingMethod != nil && r.CookingMethod.Value == other.CookingMethod.Value &&
		r.Cuisine != nil && r.Cuisine.Value == other.Cuisine.Value &&
		r.DateCreated == other.DateCreated &&
		r.DateModified == other.DateModified &&
		r.DatePublished == other.DatePublished &&
		r.Description != nil && r.Description.Value == other.Description.Value &&
		r.Keywords != nil && r.Keywords.Values == other.Keywords.Values &&
		r.Image != nil && r.Image.Value == other.Image.Value &&
		r.Ingredients != nil && slices.Equal(r.Ingredients.Values, other.Ingredients.Values) &&
		r.Instructions != nil && slices.Equal(r.Instructions.Values, other.Instructions.Values) &&
		r.Name == other.Name &&
		r.NutritionSchema != nil && r.NutritionSchema.Equal(*other.NutritionSchema) &&
		r.PrepTime == other.PrepTime &&
		r.Tools != nil && slices.Equal(r.Tools.Values, other.Tools.Values) &&
		r.TotalTime == other.TotalTime &&
		r.Yield != nil && r.Yield.Value == other.Yield.Value &&
		r.URL == other.URL
}

// NewRecipeSchema creates an initialized RecipeSchema.
func NewRecipeSchema() RecipeSchema {
	return RecipeSchema{
		AtContext:       "https://schema.org",
		AtType:          &SchemaType{Value: "Recipe"},
		Category:        NewCategory(""),
		CookingMethod:   &CookingMethod{},
		Cuisine:         &Cuisine{},
		Description:     &Description{},
		Keywords:        &Keywords{},
		Image:           &Image{},
		Ingredients:     &Ingredients{Values: make([]string, 0)},
		Instructions:    &Instructions{Values: make([]HowToItem, 0)},
		NutritionSchema: &NutritionSchema{},
		Tools:           &Tools{Values: make([]HowToItem, 0)},
		Yield:           &Yield{Value: 1},
	}
}

// HowToItem is a representation of the HowToItem schema (https://schema.org/HowToItem).
type HowToItem struct {
	Image    string `json:"image,omitempty"`
	Quantity int    `json:"requiredQuantity,omitempty"`
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Type     string `json:"@type,omitempty"`
	URL      string `json:"url,omitempty"`
}

// StringQuantity stringifies the HowToItem as `{Quantity} {Name}`.
func (h *HowToItem) StringQuantity() string {
	return strconv.Itoa(h.Quantity) + " " + h.Text
}

// NewHowToStep creates an initialized HowToStep struct.
func NewHowToStep(text string, opts ...*HowToItem) HowToItem {
	v := HowToItem{
		Type: "HowToStep",
		Text: text,
	}
	for _, opt := range opts {
		v.Name = opt.Name
		v.URL = opt.URL
		v.Image = opt.Image
	}
	return v
}

// NewHowToTool creates an initialized HowToTool struct.
func NewHowToTool(text string, opts ...*HowToItem) HowToItem {
	v := HowToItem{
		Type: "HowToTool",
		Text: text,
	}

	for _, opt := range opts {
		v.Name = opt.Name
		v.URL = opt.URL
		v.Image = opt.Image
		v.Quantity = opt.Quantity
	}

	if v.Quantity == 0 {
		v.Quantity = 1
	}

	return v
}

// Recipe transforms the RecipeSchema to a Recipe.
func (r *RecipeSchema) Recipe() (*Recipe, error) {
	if r.AtType.Value != "Recipe" {
		return nil, fmt.Errorf("RecipeSchema %#v is not based on the Schema or the field is missing", r)
	}

	var category string
	if r.Category != nil {
		if r.Category.Value == "" {
			category = "uncategorized"
		} else {
			category = strings.TrimSpace(strings.ToLower(r.Category.Value))
		}
	}

	times, err := NewTimes(r.PrepTime, r.CookTime)
	if err != nil {
		return nil, err
	}

	var nutrition Nutrition
	if r.NutritionSchema != nil {
		nutrition, err = r.NutritionSchema.nutrition()
		if err != nil {
			return nil, err
		}
	}

	created := r.DateCreated
	if r.DatePublished != "" {
		created = r.DatePublished
	}

	var createdAt time.Time
	if created != "" {
		before, _, found := strings.Cut(created, "T")
		if !found {
			before, _, _ = strings.Cut(created, " ")
		}

		createdAt, err = time.Parse(time.DateOnly, strings.ReplaceAll(before, "/", "-"))
		if err != nil {
			return nil, fmt.Errorf("could not parse createdAt date %s: %w", created, err)
		}
	}

	var images []uuid.UUID
	if r.Image != nil {
		img, err := uuid.Parse(r.Image.Value)
		if err == nil && img != uuid.Nil {
			images = append(images, img)
		}
	}

	updatedAt := createdAt
	if r.DateModified != "" {
		before, _, found := strings.Cut(r.DateModified, "T")
		if !found {
			before, _, _ = strings.Cut(r.DateModified, " ")
		}

		updatedAt, err = time.Parse(time.DateOnly, strings.ReplaceAll(before, "/", "-"))
		if err != nil {
			return nil, fmt.Errorf("could not parse modifiedAt date %q: %w", r.DateModified, err)
		}
	}

	var instructions []string
	if r.Ingredients != nil {
		instructions = make([]string, 0, len(r.Instructions.Values))
		for _, v := range r.Instructions.Values {
			instructions = append(instructions, v.Text)
		}
	}

	var keywords []string
	if r.Keywords != nil {
		keywords = extensions.Unique(strings.Split(r.Keywords.Values, ","))
	}

	var cuisine string
	if r.Cuisine != nil {
		cuisine = r.Cuisine.Value
	}

	var description string
	if r.Description != nil {
		description = r.Description.Value
	}

	var ingredients []string
	if r.Ingredients != nil {
		ingredients = r.Ingredients.Values
	}

	var tools []HowToItem
	if r.Tools != nil {
		tools = r.Tools.Values
	}

	var yield int16
	if r.Yield != nil {
		yield = r.Yield.Value
	}

	var videos []VideoObject
	if r.Video != nil {
		videos = r.Video.Values
	}

	recipe := Recipe{
		Category:     category,
		CreatedAt:    createdAt,
		Cuisine:      cuisine,
		Description:  description,
		ID:           0,
		Images:       images,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     keywords,
		Name:         r.Name,
		Nutrition:    nutrition,
		Times:        times,
		Tools:        tools,
		UpdatedAt:    updatedAt,
		URL:          r.URL,
		Videos:       videos,
		Yield:        yield,
	}

	recipe.Normalize()
	return &recipe, nil
}

// SchemaType holds the type of the schema. It should be "Recipe".
type SchemaType struct {
	Value string
}

// MarshalJSON encodes the schema's type.
func (s *SchemaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

// UnmarshalJSON decodes the type of the schema.
// The type "Recipe" will be searched for if the data is an array.
func (s *SchemaType) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		s.Value = x
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			s.Value = v.(string)
		}
	case []any:
		for _, t := range x {
			if t.(string) == "Recipe" {
				s.Value = "Recipe"
			}
		}
	}

	return nil
}

// Category holds a recipe's category.
type Category struct {
	Value string
}

// MarshalJSON encodes the category.
func (c *Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the category according to the schema (https://schema.org/recipeCategory).
// The schema specifies that the expected value is of type Text. However, some
// websites send the category in an array, which explains the need for this function.
func (c *Category) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			c.Value = v.(string)
		}
	}

	if c.Value != "" {
		c.Value = strings.ReplaceAll(c.Value, "&amp;", "&")
	}
	return nil
}

// NewCategory creates an initialized Category. The default value of a Category is 'uncategorized'.
func NewCategory(name string) *Category {
	if name == "" {
		name = "uncategorized"
	}
	return &Category{Value: name}
}

// CookingMethod holds a recipe's cooking method.
type CookingMethod struct {
	Value string
}

// MarshalJSON encodes the cuisine.
func (c *CookingMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cooking method according to the schema (https://schema.org/cookingMethod).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *CookingMethod) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	}
	return nil
}

// Cuisine holds a recipe's cuisine type.
type Cuisine struct {
	Value string
}

// MarshalJSON encodes the cuisine.
func (c *Cuisine) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Value)
}

// UnmarshalJSON decodes the cuisine according to the schema (https://schema.org/recipeCuisine).
// The schema specifies that the expected value is of type Text. However, some
// websites send the cuisine in an array, which explains the need for this function.
func (c *Cuisine) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		c.Value = x
	case []any:
		if len(x) > 0 {
			c.Value = x[0].(string)
		}
	}
	return nil
}

// Description holds a description.
type Description struct {
	Value string
}

// MarshalJSON encodes the description.
func (d *Description) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Value)
}

// UnmarshalJSON decodes the description field of a recipe.
func (d *Description) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		m := make(map[string]string)
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		v, ok := m["Value"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		s = v
	}

	s = strings.TrimSpace(s)
	replace := []struct {
		old string
		new string
	}{
		{old: "\u202f", new: ""},
		{old: "\u00a0", new: ""},
		{old: "&quot;", new: ""},
		{old: "”", new: `"`},
		{old: "\u00ad", new: ""},
	}
	for _, r := range replace {
		s = strings.ReplaceAll(s, r.old, r.new)
	}

	d.Value = s
	return nil
}

// Keywords holds keywords as a single string.
type Keywords struct {
	Values string
}

// MarshalJSON encodes the keywords.
func (k *Keywords) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Values)
}

// UnmarshalJSON decodes the recipe's keywords according to the schema (https://schema.org/keywords).
// Some websites store the keywords in an array, which explains the need
// for a custom decoder.
func (k *Keywords) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		k.Values = strings.TrimSpace(x)
	case []any:
		var xs []string
		for _, v := range x {
			xs = append(xs, v.(string))
		}
		k.Values = strings.TrimSpace(strings.Join(xs, ","))
	}

	k.Values = strings.TrimRight(k.Values, ",")
	return nil
}

// Image holds a recipe's image.
type Image struct {
	Value string
}

// MarshalJSON encodes the image.
func (i *Image) MarshalJSON() ([]byte, error) {
	s := i.Value
	if s != "" {
		s = app.Config.Address() + "/data/images/" + s
	}
	return json.Marshal(s)
}

// UnmarshalJSON decodes the image according to the schema (https://schema.org/image).
func (i *Image) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		i.Value = x
	case []any:
		if len(x) > 0 {
			switch y := x[len(x)-1].(type) {
			case string:
				i.Value = y
			case map[string]any:
				u, ok := y["url"]
				if ok {
					i.Value = u.(string)
				}
			}
		}
	case map[string]any:
		url, ok := v.(map[string]any)["url"]
		if ok {
			i.Value = url.(string)
			break
		}

		val, ok := v.(map[string]any)["Value"]
		if ok {
			i.Value = val.(string)
		}
	}
	return nil
}

// Ingredients holds a recipe's list of ingredients.
type Ingredients struct {
	Values []string
}

// MarshalJSON encodes the ingredients.
func (i *Ingredients) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

// UnmarshalJSON decodes the ingredients according to the schema (https://schema.org/recipeInstructions).
func (i *Ingredients) UnmarshalJSON(data []byte) error {
	var xv []any
	err := json.Unmarshal(data, &xv)
	if err != nil {
		m := make(map[string][]any)
		err := json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		v, ok := m["Values"]
		if !ok {
			return fmt.Errorf("could not decode description %q", data)
		}
		xv = v
	}

	cases := []struct {
		old string
		new string
	}{
		{old: "  ", new: " "},
		{old: "\u00a0", new: " "},
		{old: "&frac12;", new: "½"},
		{old: "&frac34;", new: "¾"},
		{old: "&apos;", new: "'"},
		{old: "&nbsp;", new: ""},
		{old: "&#224;", new: "à"},
		{old: "&#8217;", new: "'"},
		{old: "&#339;", new: "œ"},
		{old: "&#233;", new: "é"},
		{old: "&#239;", new: "ï"},
	}

	for _, v := range xv {
		str := v.(string)
		if str == " " {
			continue
		}

		str = strings.TrimSpace(v.(string))
		for _, c := range cases {
			str = strings.ReplaceAll(str, c.old, c.new)
		}
		i.Values = append(i.Values, str)
	}
	return nil
}

// Instructions holds a recipe's list of instructions.
type Instructions struct {
	Values []HowToItem
}

// MarshalJSON encodes the list of instructions.
func (i *Instructions) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Values)
}

// UnmarshalJSON decodes the instructions according to the schema (https://schema.org/recipeInstructions).
func (i *Instructions) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		parts := strings.Split(strings.TrimSpace(x), "\n")
		for _, s := range parts {
			if s != "" {
				i.Values = append(i.Values, NewHowToStep(strings.TrimSpace(s)))
			}
		}
	case []any:
		for _, part := range x {
			switch y := part.(type) {
			case string:
				i.Values = append(i.Values, NewHowToStep(strings.TrimSpace(y)))
			case map[string]any:
				text, ok := y["text"]
				if ok {
					str := strings.TrimSuffix(text.(string), "\n")
					i.Values = append(i.Values, NewHowToStep(strings.TrimSpace(str)))
					continue
				}

				parseSections(part, i)
			case []any:
				for _, sect := range y {
					parseSections(sect, i)
				}
			default:
				parseSections(y, i)
			}
		}
	case map[string]any:
		xv, ok := x["Values"]
		if ok {
			for _, s := range xv.([]any) {
				i.Values = append(i.Values, NewHowToStep(s.(string)))
			}
		}
	}

	cases := map[string]string{
		"\u00a0":  " ",
		"\u2009":  "",
		"&apos;":  "'",
		"&nbsp;":  " ",
		"<br>":    " ",
		"&quot;":  `"`,
		"º":       "°",
		"&#233;":  "é",
		"&#232;":  "è",
		"&#8217;": "'",
		"&#224;":  "à",
		"&#239;":  "ï",
		"&#244;":  "ô",
		"&#x27;":  "'",
	}
	for i2, value := range i.Values {
		for old, newValue := range cases {
			value.Text = strings.ReplaceAll(value.Text, old, newValue)
		}

		value.Text = strings.TrimSpace(strings.ReplaceAll(value.Text, "  ", " "))
		i.Values[i2] = value
	}

	return nil
}

type section struct {
	AtType string           `json:"@type"`
	Name   string           `json:"name"`
	Items  []map[string]any `json:"itemListElement"`
	Text   string           `json:"text"`
}

func parseSections(part any, instructions *Instructions) {
	b, err := json.Marshal(part)
	if err != nil {
		slog.Error("Marshal part failed", "error", err)
		return
	}

	var sect section
	err = json.Unmarshal(b, &sect)
	if err != nil {
		slog.Error("Marshal section failed", "error", err)
		return
	}

	if sect.Text != "" {
		instructions.Values = append(instructions.Values, NewHowToStep(sect.Text))
	}

	for _, item := range sect.Items {
		text, ok := item["text"]
		if ok {
			var str string
			switch x := text.(type) {
			case string:
				str = x
			case []any:
				xs := make([]string, 0, len(x))
				for _, v := range x {
					xs = append(xs, v.(string))
				}
				str = strings.Join(xs, ", ")
			default:
				continue
			}

			str = strings.TrimSuffix(str, "\n")
			str = strings.ReplaceAll(str, "\u00a0", "")
			str = strings.ReplaceAll(str, "\u2009", "")
			instructions.Values = append(instructions.Values, NewHowToStep(strings.TrimSpace(str)))
		}
	}
}

// Yield holds a recipe's yield.
type Yield struct {
	Value int16
}

// MarshalJSON encodes the value of the yield.
func (y *Yield) MarshalJSON() ([]byte, error) {
	return json.Marshal(y.Value)
}

// UnmarshalJSON decodes the yield according to the schema (https://schema.org/recipeYield).
func (y *Yield) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		parts := strings.Split(v.(string), " ")
		for _, part := range parts {
			i, err := strconv.ParseInt(part, 10, 8)
			if err == nil {
				y.Value = int16(i)
				break
			}
		}
	case float64:
		y.Value = int16(x)
	case []any:
		if len(x) == 0 {
			break
		}

		switch t := x[0].(type) {
		case float64:
			if t >= math.MinInt16 && t <= math.MaxInt16 {
				y.Value = int16(t)
			} else {
				y.Value = int16(4)
			}
		case string:
			parts := strings.Split(t, " ")
			for _, part := range parts {
				i, err := strconv.ParseInt(part, 10, 8)
				if err == nil {
					y.Value = int16(i)
					break
				}
			}
		}
	case map[string]any:
		v, ok := x["Value"]
		if ok {
			y.Value = int16(v.(float64))
		}
	}
	return nil
}

// Videos holds a list of VideoObject.
type Videos struct {
	Values []VideoObject
}

// MarshalJSON encodes the value of the yield.
func (v *Videos) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Values)
}

// UnmarshalJSON decodes the Videos according to the schema (https://schema.org/VideoObject).
func (v *Videos) UnmarshalJSON(data []byte) error {
	var t any
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	switch x := t.(type) {
	case map[string]any:
		tm, ok := x["uploadDate"]
		if ok {
			parsed, _ := parseTime(tm.(string))
			x["uploadDate"] = parsed.Format(time.RFC3339)
			data, _ = json.Marshal(x)
		}

		var vid VideoObject
		err = json.Unmarshal(data, &vid)
		if err != nil {
			return err
		}
		v.Values = append(v.Values, vid)
	case []any:
		for _, a := range x {
			tm, ok := a.(map[string]any)["uploadDate"]
			if ok {
				parsed, _ := parseTime(tm.(string))
				a.(map[string]any)["uploadDate"] = parsed.Format(time.RFC3339)
				data, _ = json.Marshal(x)
			}

			xb, _ := json.Marshal(a)
			var vid VideoObject
			err = json.Unmarshal(xb, &vid)
			if err != nil {
				return err
			}
			v.Values = append(v.Values, vid)
		}
	}
	return nil
}

func parseTime(tm string) (time.Time, error) {
	var (
		parsed time.Time
		err    error
	)

	layouts := []string{"01-02-2006", time.RFC3339, "2006-01-02T15:04:05.999", "2006/01/02"}
	for _, layout := range layouts {
		parsed, err = time.Parse(layout, tm)
		if err == nil {
			break
		}
	}

	return parsed, err
}

// VideoObject is a representation of the VideoObject schema (https://schema.org/VideoObject).
type VideoObject struct {
	AtType       string        `json:"@type"`
	ContentURL   string        `json:"contentUrl,omitempty"`
	Description  string        `json:"description,omitempty"`
	Duration     string        `json:"duration,omitempty"`
	EmbedURL     string        `json:"embedUrl,omitempty"`
	Expires      time.Time     `json:"expires,omitempty"`
	ID           uuid.UUID     `json:"-"`
	IsIFrame     bool          `json:"_"`
	Name         string        `json:"name,omitempty"`
	ThumbnailURL *ThumbnailURL `json:"thumbnailUrl,omitempty"`
	UploadDate   time.Time     `json:"uploadDate,omitempty"`
}

// NutritionSchema is a representation of the nutrition schema (https://schema.org/NutritionInformation).
type NutritionSchema struct {
	Calories       string `json:"calories,omitempty"`
	Carbohydrates  string `json:"carbohydrateContent,omitempty"`
	Cholesterol    string `json:"cholesterolContent,omitempty"`
	Fat            string `json:"fatContent,omitempty"`
	Fiber          string `json:"fiberContent,omitempty"`
	Protein        string `json:"proteinContent,omitempty"`
	SaturatedFat   string `json:"saturatedFatContent,omitempty"`
	Servings       string `json:"servingSize,omitempty"`
	Sodium         string `json:"sodiumContent,omitempty"`
	Sugar          string `json:"sugarContent,omitempty"`
	TransFat       string `json:"transFatContent,omitempty"`
	UnsaturatedFat string `json:"unsaturatedFatContent,omitempty"`
}

// Equal verifies whether the NutritionSchema equals the other.
func (n *NutritionSchema) Equal(other NutritionSchema) bool {
	return n.Calories == other.Calories &&
		n.Carbohydrates == other.Carbohydrates &&
		n.Cholesterol == other.Cholesterol &&
		n.Fat == other.Fat &&
		n.Fiber == other.Fiber &&
		n.Protein == other.Protein &&
		n.SaturatedFat == other.SaturatedFat &&
		n.Servings == other.Servings &&
		n.Sodium == other.Sodium &&
		n.Sugar == other.Sugar &&
		n.TransFat == other.TransFat &&
		n.UnsaturatedFat == other.UnsaturatedFat
}

func (n *NutritionSchema) nutrition() (Nutrition, error) {
	nutrition := Nutrition{
		Calories:           n.Calories,
		Cholesterol:        n.Cholesterol,
		Fiber:              n.Fiber,
		Protein:            n.Protein,
		SaturatedFat:       n.SaturatedFat,
		Sodium:             n.Sodium,
		Sugars:             n.Sugar,
		TotalCarbohydrates: n.Carbohydrates,
		TotalFat:           n.Fat,
		UnsaturatedFat:     n.UnsaturatedFat,
	}

	if n.Servings != "" {
		servings, err := strconv.ParseInt(n.Servings, 10, 16)
		if err != nil {
			return Nutrition{}, nil
		}

		nutrition.IsPerServing = true
		nutrition.Scale(1 / float64(servings))
	}

	return nutrition, nil
}

// UnmarshalJSON decodes the nutrition according to the schema.
func (n *NutritionSchema) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case map[string]any:
		if val, ok := x["calories"].(string); ok {
			n.Calories = strings.TrimSpace(val)
		}

		if val, ok := x["carbohydrateContent"].(string); ok {
			n.Carbohydrates = val
		}

		if val, ok := x["cholesterolContent"].(string); ok {
			n.Cholesterol = val
		}

		if val, ok := x["fatContent"].(string); ok {
			n.Fat = val
		}

		if val, ok := x["fiberContent"].(string); ok {
			n.Fiber = val
		}

		if val, ok := x["proteinContent"].(string); ok {
			n.Protein = val
		}

		if val, ok := x["saturatedFatContent"].(string); ok {
			n.SaturatedFat = val
		}

		if val, ok := x["servingSize"].(string); ok {
			xs := strings.Split(val, " ")
			if len(xs) == 2 && len(xs[1]) < 4 {
				n.Servings = val
			} else {
				for _, s := range xs {
					_, err := strconv.Atoi(s)
					if err == nil {
						n.Servings = s
						break
					}
				}
			}
		}

		if val, ok := x["sodiumContent"].(string); ok {
			n.Sodium = val
		}

		if val, ok := x["sugarContent"].(string); ok {
			n.Sugar = val
		}

		if val, ok := x["transFatContent"].(string); ok {
			n.TransFat = val
		}

		if val, ok := x["unsaturatedFatContent"].(string); ok {
			n.UnsaturatedFat = val
		}
	default:
		slog.Warn("Could not parse nutrition schema", "schema", v, "type", x)
	}
	return nil
}

// ThumbnailURL holds a recipe's thumbnail.
type ThumbnailURL struct {
	Value string
}

// MarshalJSON encodes the thumbnail URL.
func (t *ThumbnailURL) MarshalJSON() ([]byte, error) {
	s := t.Value
	if s != "" {
		s = app.Config.Address() + "/data/images/thumbnails/" + s
	}
	return json.Marshal(s)
}

// UnmarshalJSON decodes the thumbnail URL according to the schema (https://schema.org/URL).
func (t *ThumbnailURL) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		t.Value = x
	case []any:
		if len(x) > 0 {
			t.Value = x[0].(string)
		}
	}
	return nil
}

// Tools holds the list of tools used for a recipe.
type Tools struct {
	Values []HowToItem
}

// MarshalJSON encodes the list of tools.
func (t *Tools) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Values)
}

// UnmarshalJSON decodes the tools according to the schema (https://schema.org/tool).
func (t *Tools) UnmarshalJSON(data []byte) error {
	var v any
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	switch x := v.(type) {
	case string:
		delim := ","
		if strings.Contains(x, "\n") {
			delim = "\n"
		}

		parts := strings.Split(x, delim)
		for _, s := range parts {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}

			var q int
			subParts := strings.Split(s, " ")
			if len(subParts) > 1 {
				parsed, err := strconv.Atoi(subParts[0])
				if err == nil {
					q = parsed
					s = strings.Join(subParts[1:], " ")
				}
			}

			t.Values = append(t.Values, NewHowToTool(s, &HowToItem{Quantity: q}))
		}
	case []any:
		for _, a := range x {
			switch x2 := a.(type) {
			case map[string]any:
				xb, err := json.Marshal(a)
				if err != nil {
					return err
				}

				var tool HowToItem
				err = json.Unmarshal(xb, &tool)
				if err != nil {
					return err
				}

				if tool.Text == "" {
					item, ok := x2["item"]
					if ok {
						tool.Text = item.(string)
					}
				}

				t.Values = append(t.Values, tool)
			}
		}
	case map[string]any:
		var tool HowToItem
		err = json.Unmarshal(data, &tool)
		if err != nil {
			return err
		}
		t.Values = append(t.Values, tool)
	default:
		slog.Warn("Could not parse Tool schema", "data", data)
	}

	return nil
}
