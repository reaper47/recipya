package scraper

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reaper47/recipya/internal/models"
)

type bettibossi struct {
	RezeptKopf struct {
		ErstellungsDatum string `json:"ErstellungsDatum"`
		MutationsDatum   string `json:"MutationsDatum"`
		LiveDatum        string `json:"LiveDatum"`
		Sprache          string `json:"Sprache"`
		Titel            string `json:"Titel"`
		Beschreibung     string `json:"Beschreibung"`
		Menge1           int    `json:"Menge_1"`
		Menge1Einheit    string `json:"Menge_1_Einheit"`
		Menge2           int    `json:"Menge_2"`
		Menge2Einheit    string `json:"Menge_2_Einheit"`
		SkalierenNach    int    `json:"skalieren_nach"`
	} `json:"RezeptKopf"`
	Zeiten []struct {
		Wert       int    `json:"Wert"`
		Einheit    string `json:"Einheit"`
		Kategorie  string `json:"Kategorie"`
		Zeitangabe string `json:"Zeitangabe"`
	} `json:"Zeiten"`
	Naehrwerte []struct {
		Kurzbezeichnung string `json:"Kurzbezeichnung"`
		Bezeichnung     string `json:"Bezeichnung"`
		Wert1           int    `json:"Wert_1"`
		Wert1Einheit    string `json:"Wert_1_Einheit"`
	} `json:"Naehrwerte"`
	Kategorisierungen []struct {
		GrpBez string `json:"GrpBez"`
		Bez    string `json:"Bez"`
	} `json:"Kategorisierungen"`
	Subrezepte []struct {
		Schritte []struct {
			Zutaten []struct {
				Menge   float64 `json:"Menge"`
				Einheit string  `json:"Einheit"`
				Zutat   string  `json:"Zutat"`
			} `json:"Zutaten,omitempty"`
			Anleitungen []struct {
				Text string `json:"Text"`
			} `json:"Anleitungen"`
		} `json:"Schritte"`
	} `json:"Subrezepte"`
}

func scrapeBettybossi(root *goquery.Document) (models.RecipeSchema, error) {
	var (
		b   bettibossi
		err error
		rs  = models.NewRecipeSchema()
	)

	rs.Image.Value = getPropertyContent(root, "og:image")

	root.Find("meta").Each(func(_ int, s *goquery.Selection) {
		n, ok := s.Attr("data-rjson")
		if ok {
			err = json.Unmarshal([]byte(n), &b)
			if err != nil {
				return
			}
		}
	})

	var (
		prepDur time.Duration
		cookDur time.Duration
	)
	for _, z := range b.Zeiten {
		if prepDur != 0 && cookDur != 0 {
			break
		} else if z.Wert == 0 {
			continue
		}

		lower := strings.ToLower(z.Kategorie)
		if strings.Contains(lower, "zubereiten") {
			if z.Einheit == "Min." {
				prepDur += time.Duration(z.Wert) * time.Minute
			} else {
				prepDur += time.Duration(z.Wert) * time.Hour
			}
		} else if strings.Contains(lower, "backen") {
			if z.Einheit == "Min." {
				cookDur += time.Duration(z.Wert) * time.Minute
			} else {
				cookDur += time.Duration(z.Wert) * time.Hour
			}
		}
	}

	if prepDur > 0 {
		prep := "PT" + strings.Replace(prepDur.String(), "h", "H", 1)
		prep = strings.Replace(prep, "m", "M", 1)
		before, _, ok := strings.Cut(prep, "M")
		if ok {
			prep = before + "M"
		}
		rs.PrepTime = prep
	}

	if cookDur > 0 {
		cook := "PT" + strings.Replace(cookDur.String(), "h", "H", 1)
		cook = strings.Replace(cook, "m", "M", 1)
		before, _, ok := strings.Cut(cook, "M")
		if ok {
			cook = before + "M"
		}
		rs.CookTime = cook
	}

	var keywords []string
	for _, k := range b.Kategorisierungen {
		if k.GrpBez == "Eigenschaft" || k.GrpBez == "Saison" {
			keywords = append(keywords, k.Bez)
		} else if k.Bez == "Gericht / Gang" && rs.Category.Value == "" {
			rs.Category.Value = k.Bez
		}
	}

	var ns models.NutritionSchema
	for _, n := range b.Naehrwerte {
		v := strconv.Itoa(n.Wert1) // + " " + n.Wert1Einheit 	2024-07-01: Commented out this section (left the original code as it was in German in case we need to revert) as we only scrape the digits, as part of https://github.com/reaper47/recipya/pull/382
		switch n.Kurzbezeichnung {
		case "E":
			if n.Bezeichnung == "Eiweiss" {
				ns.Protein = v
			} else {
				ns.Calories = v
			}
		case "F":
			ns.Fat = v
		case "Kh":
			ns.Carbohydrates = v
		case "su":
			ns.Sugar = v
		case "FATS":
			ns.SaturatedFat = v
		case "fib":
			ns.Fiber = v
		case "NaCl":
			ns.Sodium = v
		}
	}

	if len(b.Subrezepte) > 0 {
		for _, s := range b.Subrezepte[0].Schritte {
			for _, ins := range s.Anleitungen {
				rs.Instructions.Values = append(rs.Instructions.Values, models.NewHowToStep(ins.Text))
			}

			for _, ing := range s.Zutaten {
				rs.Ingredients.Values = append(rs.Ingredients.Values, strconv.Itoa(int(ing.Menge))+" "+ing.Einheit+" "+ing.Zutat)
			}
		}
	}

	rs.DateCreated = b.RezeptKopf.ErstellungsDatum
	rs.DateModified = b.RezeptKopf.LiveDatum
	rs.DatePublished = b.RezeptKopf.LiveDatum
	rs.Description = &models.Description{Value: b.RezeptKopf.Beschreibung}
	rs.Keywords = &models.Keywords{Values: strings.Join(keywords, ",")}
	rs.Name = b.RezeptKopf.Titel
	rs.NutritionSchema = &ns
	rs.Yield = &models.Yield{Value: int16(b.RezeptKopf.Menge1)}

	return rs, err
}
