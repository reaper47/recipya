package units_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/units"
	"math"
	"testing"
)

func TestConvertSentence(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		from units.System
		to   units.System
		want string
	}{
		{
			name: "celsius to fahrenheit",
			in:   "Rope with 200 degrees celsius of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 392 °F of rice.",
		},
		{
			name: "fahrenheit to celsius",
			in:   "Rope with 475 degrees fahrenheit of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 246 °C of rice.",
		},
		{
			name: "mm to yards",
			in:   "Rope with 455980 mm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 498.67 yards of rice.",
		},
		{
			name: "mm to feet",
			in:   "Rope with 30800mm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 101.05 feet of rice.",
		},
		{
			name: "mm to inches",
			in:   "Rope with 666 mm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 26.22 inches of rice.",
		},
		{
			name: "cm to yards",
			in:   "Rope with 45598 cm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 498.67 yards of rice.",
		},
		{
			name: "cm to feet",
			in:   "Rope with 3080cm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 101.05 feet of rice.",
		},
		{
			name: "cm to inches",
			in:   "Rope with 66 cm of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 25.98 inches of rice.",
		},
		{
			name: "m to yards",
			in:   "Rope with 450 m of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 492.13 yards of rice.",
		},
		{
			name: "m to feet",
			in:   "Rope with 300m of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 984.25 feet of rice.",
		},
		{
			name: "m to inches",
			in:   "Rope with 6 M of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 236.22 inches of rice.",
		},
		{
			name: "inches to mm",
			in:   "Rope with 0.2 in of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 5.08 mm of rice.",
		},
		{
			name: "inches to cm",
			in:   "Rope with 23 in of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 58.42 cm of rice.",
		},
		{
			name: "inches to m",
			in:   "Rope with 600 in of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 15.24 m of rice.",
		},
		{
			name: "feet to mm",
			in:   "Rope with 0.5 foot of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 15240 mm of rice.",
		},
		{
			name: "feet to cm",
			in:   "Rope with 48′ of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 1463.04 cm of rice.",
		},
		{
			name: "feet to m",
			in:   "Rope with 600ft of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 182.88 m of rice.",
		},
		{
			name: "yard to mm",
			in:   "Rope with 0.1 yards of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 91.44 mm of rice.",
		},
		{
			name: "yard to cm",
			in:   "Rope with 1yard of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 91.44 cm of rice.",
		},
		{
			name: "yard to m",
			in:   "Rope with 2 yard of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 1.83 m of rice.",
		},
		{
			name: "mg to lb",
			in:   "Rope with 20000 mg of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 0.04 lb of rice.",
		},
		{
			name: "mg to oz",
			in:   "Rope with 48349.52 milligrams of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.71 oz of rice.",
		},
		{
			name: "g to lb",
			in:   "Rope with 200 g of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 0.44 lb of rice.",
		},
		{
			name: "g to oz",
			in:   "Rope with 4834 grams of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 170.51 oz of rice.",
		},
		{
			name: "kg to lb",
			in:   "Rope with 2 kg of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 4.41 lb of rice.",
		},
		{
			name: "kg to oz",
			in:   "Rope with 44 kilograms of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1552.05 oz of rice.",
		},
		{
			name: "lb to mg",
			in:   "Rope with 0.0004 lb of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 181.44 mg of rice.",
		},
		{
			name: "lb to g",
			in:   "Rope with 1.5 lb of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 680.39 g of rice.",
		},
		{
			name: "lb to kg",
			in:   "Rope with 44 lb of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 19.96 kg of rice.",
		},
		{
			name: "oz to mg",
			in:   "Rope with 0.0004 oz of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 11.34 mg of rice.",
		},
		{
			name: "oz to g",
			in:   "Rope with 2 oz of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 56.7 g of rice.",
		},
		{
			name: "oz to kg",
			in:   "Rope with 44 oz of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 1.25 kg of rice.",
		},
		{
			name: "dL to tsp",
			in:   "Rope with 0.1dL of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 2.03 tsp of rice.",
		},
		{
			name: "dL to tbsp",
			in:   "Rope with 0.25dl of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 tbsp of rice.",
		},
		{
			name: "dL to fl oz",
			in:   "Rope with 1 decilitres of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 3.38 fl oz of rice.",
		},
		{
			name: "dL to cup",
			in:   "Rope with 4dl of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 cups of rice.",
		},
		{
			name: "dL to pt",
			in:   "Rope with 7dl of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.48 pints of rice.",
		},
		{
			name: "dL to qt",
			in:   "Rope with 30dl of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 3.17 fl qt of rice.",
		},
		{
			name: "dL to gal",
			in:   "Rope with 40 dL of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 10.57 gallons of rice.",
		},
		{
			name: "mL to tsp",
			in:   "Rope with 5mL of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1 tsp of rice.",
		},
		{
			name: "mL to tbsp",
			in:   "Rope with 25ml of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.67 tbsp of rice.",
		},
		{
			name: "mL to fl oz",
			in:   "Rope with 50 millilitres of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 fl oz of rice.",
		},
		{
			name: "mL to cup",
			in:   "Rope with 800ml of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 3.38 cups of rice.",
		},
		{
			name: "mL to pt",
			in:   "Rope with 1500ml of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 3.17 pints of rice.",
		},
		{
			name: "mL to qt",
			in:   "Rope with 2500ml of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 2.64 fl qt of rice.",
		},
		{
			name: "mL to gal",
			in:   "Rope with 40000 mL of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 10.57 gallons of rice.",
		},
		{
			name: "l to tsp",
			in:   "Rope with 0.005L of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.01 tsp of rice.",
		},
		{
			name: "l to tbsp",
			in:   "Rope with 0.02l of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.35 tbsp of rice.",
		},
		{
			name: "l to fl oz",
			in:   "Rope with 0.05 litres of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 fl oz of rice.",
		},
		{
			name: "l to cup",
			in:   "Rope with 0.4l of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 cups of rice.",
		},
		{
			name: "l to pt",
			in:   "Rope with 0.8l of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.69 pints of rice.",
		},
		{
			name: "l to qt",
			in:   "Rope with 3l of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 3.17 fl qt of rice.",
		},
		{
			name: "l to gal",
			in:   "Rope with 6 L of rice.",
			from: units.MetricSystem,
			to:   units.ImperialSystem,
			want: "Rope with 1.59 gallons of rice.",
		},
		{
			name: "tsp to mL",
			in:   "Rope with 1 tsp of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 5 ml of rice.",
		},
		{
			name: "tsp to dL",
			in:   "Rope with 200tsp. of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 9.86 dl of rice.",
		},
		{
			name: "tsp to L",
			in:   "Rope with 250 teaspoons of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 1.23 L of rice.",
		},
		{
			name: "tbsp to mL",
			in:   "Rope with 1 tbsp of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 14.79 ml of rice.",
		},
		{
			name: "tbsp to dL",
			in:   "Rope with 10 tablespoons of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 1.48 dl of rice.",
		},
		{
			name: "tbsp to L",
			in:   "Rope with 250 tablespoons of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 3.7 L of rice.",
		},
		{
			name: "fl oz to mL",
			in:   "Rope with 0.1 fl oz of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 2.96 ml of rice.",
		},
		{
			name: "fl oz to dL",
			in:   "Rope with 1floz. of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 0.3 dL of rice.",
		},
		{
			name: "fl oz to L",
			in:   "Rope with 250 fluid ounces of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 7.39 L of rice.",
		},
		{
			name: "cup to mL",
			in:   "Rope with 0.1 cup of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 23.66 ml of rice.",
		},
		{
			name: "cup to dL",
			in:   "Rope with 1cups of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 2.37 dl of rice.",
		},
		{
			name: "cup to L",
			in:   "Rope with 250cup of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 59.15 L of rice.",
		},
		{
			name: "pt to mL",
			in:   "Rope with 0.1 pint of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 47.32 ml of rice.",
		},
		{
			name: "pt to dL",
			in:   "Rope with 1pt of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 4.73 dl of rice.",
		},
		{
			name: "pt to L",
			in:   "Rope with 250fl pt of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 118.29 L of rice.",
		},
		{
			name: "qt to mL",
			in:   "Rope with 0.1 quart of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 94.64 ml of rice.",
		},
		{
			name: "qt to dL",
			in:   "Rope with 1qt of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 9.46 dl of rice.",
		},
		{
			name: "qt to L",
			in:   "Rope with 250fl qt of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 236.59 L of rice.",
		},
		{
			name: "gallon to mL",
			in:   "Rope with 0.01 gal of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 37.85 ml of rice.",
		},
		{
			name: "gallon to dL",
			in:   "Rope with 0.1gal of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 3.79 dl of rice.",
		},
		{
			name: "gallon to L",
			in:   "Rope with 250gallons of rice.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Rope with 946.35 L of rice.",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := units.ConvertSentence(tc.in, tc.from, tc.to)
			assertEqual(t, got, tc.want)
		})
	}

	testcases2 := []struct {
		name string
		in   string
		from units.System
		to   units.System
		want string
	}{
		{
			name: "irregular",
			in:   "2 to 3 pounds chicken wings",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "0.91 to 1.36 kg chicken wings",
		},
		{
			name: "irregular",
			in:   "5 1/3 tablespoons (2/3 stick) salted butter",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "78.86 ml (2/3 stick) salted butter",
		},
		{
			name: "irregular",
			in:   "Two 12-ounce bottles Frank's Red Hot Sauce (accept NO imitations!!!!)",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Two 340.19 g bottles Frank's Red Hot Sauce (accept NO imitations!!!!)",
		},
		{
			name: "irregular",
			in:   "1-3/4 cups water",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "4.14 dl water",
		},
		{
			name: "irregular",
			in:   "One 3-pound bag frozen hash browns",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "One 1.36 kg bag frozen hash browns",
		},
		{
			name: "irregular",
			in:   "1/3 cup butter",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "78.86 ml butter",
		},
		{
			name: "irregular",
			in:   "1/4-1/2 teaspoon salt plus more for seasoning",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "3.75 ml salt plus more for seasoning",
		},
		{
			name: "irregular",
			in:   "½ teaspoon paprika",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "2.5 ml paprika",
		},
		{
			name: "irregular",
			in:   "Scrub potatoes (do not peel them). Dice into 1” cubes.",
			from: units.ImperialSystem,
			to:   units.MetricSystem,
			want: "Scrub potatoes (do not peel them). Dice into 2.54 cm cubes.",
		},
	}
	for _, tc := range testcases2 {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := units.ConvertSentence(tc.in, tc.from, tc.to)
			assertEqual(t, got, tc.want)
		})
	}
}

func TestDetectMeasurementSystemFromSentence(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		want units.System
	}{
		{
			name: "imperial",
			in:   "2 teaspoons hot water",
			want: units.ImperialSystem,
		},
		{
			name: "metric",
			in:   "2 mL hot water",
			want: units.MetricSystem,
		},
		{
			name: "invalid",
			in:   "2 oranges",
			want: units.InvalidSystem,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := units.DetectMeasurementSystemFromSentence(tc.in)
			assertEqual(t, got, tc.want)
		})
	}
}

func TestNewMeasurement(t *testing.T) {
	testcases := []struct {
		quantity int
		unit     string
		want     units.Measurement
	}{
		{quantity: 1, unit: "ml", want: units.Measurement{Quantity: 1, Unit: units.Millilitre}},
		{quantity: 1, unit: "mL", want: units.Measurement{Quantity: 1, Unit: units.Millilitre}},
		{quantity: 1, unit: "millilitres", want: units.Measurement{Quantity: 1, Unit: units.Millilitre}},
		{quantity: 1, unit: "milliliters", want: units.Measurement{Quantity: 1, Unit: units.Millilitre}},
		{quantity: 1, unit: "cc", want: units.Measurement{Quantity: 1, Unit: units.Millilitre}},

		{quantity: 2, unit: "l", want: units.Measurement{Quantity: 2, Unit: units.Litre}},
		{quantity: 2, unit: "L", want: units.Measurement{Quantity: 2, Unit: units.Litre}},
		{quantity: 2, unit: "liters", want: units.Measurement{Quantity: 2, Unit: units.Litre}},
		{quantity: 2, unit: "litres", want: units.Measurement{Quantity: 2, Unit: units.Litre}},

		{quantity: 3, unit: "dl", want: units.Measurement{Quantity: 3, Unit: units.Decilitre}},
		{quantity: 3, unit: "dL", want: units.Measurement{Quantity: 3, Unit: units.Decilitre}},
		{quantity: 3, unit: "deciliters", want: units.Measurement{Quantity: 3, Unit: units.Decilitre}},
		{quantity: 3, unit: "decilitres", want: units.Measurement{Quantity: 3, Unit: units.Decilitre}},

		{quantity: 4, unit: "teaspoons", want: units.Measurement{Quantity: 4, Unit: units.Teaspoon}},
		{quantity: 4, unit: "tsp", want: units.Measurement{Quantity: 4, Unit: units.Teaspoon}},
		{quantity: 4, unit: "tsp.", want: units.Measurement{Quantity: 4, Unit: units.Teaspoon}},

		{quantity: 5, unit: "tablespoons", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbl", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbl.", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tb", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbs", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbs.", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbsp", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},
		{quantity: 5, unit: "tbsp.", want: units.Measurement{Quantity: 5, Unit: units.Tablespoon}},

		{quantity: 6, unit: "fluid ounces", want: units.Measurement{Quantity: 6, Unit: units.FlOz}},
		{quantity: 6, unit: "fl oz", want: units.Measurement{Quantity: 6, Unit: units.FlOz}},
		{quantity: 6, unit: "fl. oz.", want: units.Measurement{Quantity: 6, Unit: units.FlOz}},

		{quantity: 8, unit: "cups", want: units.Measurement{Quantity: 8, Unit: units.Cup}},

		{quantity: 9, unit: "pints", want: units.Measurement{Quantity: 9, Unit: units.Pint}},
		{quantity: 9, unit: "pt", want: units.Measurement{Quantity: 9, Unit: units.Pint}},
		{quantity: 9, unit: "fl pt", want: units.Measurement{Quantity: 9, Unit: units.Pint}},
		{quantity: 9, unit: "fl. pt.", want: units.Measurement{Quantity: 9, Unit: units.Pint}},

		{quantity: 10, unit: "quarts", want: units.Measurement{Quantity: 10, Unit: units.Quart}},
		{quantity: 10, unit: "qt", want: units.Measurement{Quantity: 10, Unit: units.Quart}},
		{quantity: 10, unit: "fl qt", want: units.Measurement{Quantity: 10, Unit: units.Quart}},
		{quantity: 10, unit: "fl. qt.", want: units.Measurement{Quantity: 10, Unit: units.Quart}},

		{quantity: 11, unit: "gallons", want: units.Measurement{Quantity: 11, Unit: units.Gallon}},
		{quantity: 11, unit: "gals", want: units.Measurement{Quantity: 11, Unit: units.Gallon}},

		{quantity: 12, unit: "mg", want: units.Measurement{Quantity: 12, Unit: units.Milligram}},
		{quantity: 12, unit: "milligrams", want: units.Measurement{Quantity: 12, Unit: units.Milligram}},
		{quantity: 12, unit: "milligrammes", want: units.Measurement{Quantity: 12, Unit: units.Milligram}},

		{quantity: 13, unit: "g", want: units.Measurement{Quantity: 13, Unit: units.Gram}},
		{quantity: 13, unit: "grams", want: units.Measurement{Quantity: 13, Unit: units.Gram}},
		{quantity: 13, unit: "grammes", want: units.Measurement{Quantity: 13, Unit: units.Gram}},

		{quantity: 14, unit: "kg", want: units.Measurement{Quantity: 14, Unit: units.Kilogram}},
		{quantity: 14, unit: "kilograms", want: units.Measurement{Quantity: 14, Unit: units.Kilogram}},
		{quantity: 14, unit: "kilogrammes", want: units.Measurement{Quantity: 14, Unit: units.Kilogram}},

		{quantity: 15, unit: "lb", want: units.Measurement{Quantity: 15, Unit: units.Pound}},
		{quantity: 15, unit: "#", want: units.Measurement{Quantity: 15, Unit: units.Pound}},
		{quantity: 15, unit: "pounds", want: units.Measurement{Quantity: 15, Unit: units.Pound}},

		{quantity: 16, unit: "ounce", want: units.Measurement{Quantity: 16, Unit: units.Ounce}},
		{quantity: 16, unit: "oz", want: units.Measurement{Quantity: 16, Unit: units.Ounce}},
		{quantity: 16, unit: "oz.", want: units.Measurement{Quantity: 16, Unit: units.Ounce}},

		{quantity: 17, unit: "mm", want: units.Measurement{Quantity: 17, Unit: units.Millimeter}},
		{quantity: 17, unit: "millimeters", want: units.Measurement{Quantity: 17, Unit: units.Millimeter}},
		{quantity: 17, unit: "millimetres", want: units.Measurement{Quantity: 17, Unit: units.Millimeter}},

		{quantity: 18, unit: "cm", want: units.Measurement{Quantity: 18, Unit: units.Centimeter}},
		{quantity: 18, unit: "centimeters", want: units.Measurement{Quantity: 18, Unit: units.Centimeter}},
		{quantity: 18, unit: "centimetres", want: units.Measurement{Quantity: 18, Unit: units.Centimeter}},

		{quantity: 19, unit: "m", want: units.Measurement{Quantity: 19, Unit: units.Meter}},
		{quantity: 19, unit: "meters", want: units.Measurement{Quantity: 19, Unit: units.Meter}},
		{quantity: 19, unit: "metres", want: units.Measurement{Quantity: 19, Unit: units.Meter}},

		{quantity: 20, unit: "inches", want: units.Measurement{Quantity: 20, Unit: units.Inch}},
		{quantity: 20, unit: "inch", want: units.Measurement{Quantity: 20, Unit: units.Inch}},
		{quantity: 20, unit: "in", want: units.Measurement{Quantity: 20, Unit: units.Inch}},
		{quantity: 20, unit: `"`, want: units.Measurement{Quantity: 20, Unit: units.Inch}},

		{quantity: 21, unit: "yards", want: units.Measurement{Quantity: 21, Unit: units.Yard}},

		{quantity: 22, unit: "°C", want: units.Measurement{Quantity: 22, Unit: units.Celsius}},
		{quantity: 22, unit: "degrees celsius", want: units.Measurement{Quantity: 22, Unit: units.Celsius}},
		{quantity: 22, unit: "degree celsius", want: units.Measurement{Quantity: 22, Unit: units.Celsius}},

		{quantity: 23, unit: "°F", want: units.Measurement{Quantity: 23, Unit: units.Fahrenheit}},
		{quantity: 23, unit: "degrees Farenheit", want: units.Measurement{Quantity: 23, Unit: units.Fahrenheit}},
		{quantity: 23, unit: "degree Farenheit", want: units.Measurement{Quantity: 23, Unit: units.Fahrenheit}},
		{quantity: 23, unit: "degree Fahrenheit", want: units.Measurement{Quantity: 23, Unit: units.Fahrenheit}},
		{quantity: 23, unit: "degrees Fahrenheit", want: units.Measurement{Quantity: 23, Unit: units.Fahrenheit}},

		{quantity: 24, unit: "feet", want: units.Measurement{Quantity: 24, Unit: units.Feet}},
		{quantity: 24, unit: "foot", want: units.Measurement{Quantity: 24, Unit: units.Feet}},
		{quantity: 24, unit: "ft", want: units.Measurement{Quantity: 24, Unit: units.Feet}},
		{quantity: 24, unit: "′", want: units.Measurement{Quantity: 24, Unit: units.Feet}},
	}
	for _, tc := range testcases {
		t.Run(tc.unit, func(t *testing.T) {
			got, err := units.NewMeasurement(float64(tc.quantity), tc.unit)
			if err != nil {
				wantErr := errors.New("unit " + tc.unit + " is unsupported")
				if errors.Is(err, wantErr) {
					t.Fatalf("got error %q but want %q", err, wantErr)
				}
				t.Fatalf("got unexpected error %q", err)
			}

			if !cmp.Equal(got, tc.want) {
				t.Errorf("got %#v but want %#v", got, tc.want)
			}
		})
	}
}

func TestMeasurement_Convert(t *testing.T) {
	testcases := []struct {
		name string
		to   units.Unit
		in   units.Measurement
		want units.Measurement
	}{
		{
			name: "celsius to fahrenheit",
			in:   units.Measurement{Quantity: 162, Unit: units.Celsius},
			want: units.Measurement{Quantity: 324, Unit: units.Fahrenheit},
		},
		{
			name: "celsius to celsius",
			in:   units.Measurement{Quantity: 162, Unit: units.Celsius},
			want: units.Measurement{Quantity: 162, Unit: units.Celsius},
		},
		{
			name: "fahrenheit to celsius",
			in:   units.Measurement{Quantity: 527, Unit: units.Fahrenheit},
			want: units.Measurement{Quantity: 275, Unit: units.Celsius},
		},
		{
			name: "fahrenheit to fahrenheit",
			in:   units.Measurement{Quantity: 527, Unit: units.Fahrenheit},
			want: units.Measurement{Quantity: 527, Unit: units.Fahrenheit},
		},
		{
			name: "mm to mm",
			in:   units.Measurement{Quantity: 1250, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 1250, Unit: units.Millimeter},
		},
		{
			name: "mm to cm",
			in:   units.Measurement{Quantity: 125, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 12.5, Unit: units.Centimeter},
		},
		{
			name: "mm to m",
			in:   units.Measurement{Quantity: 125, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 0.125, Unit: units.Meter},
		},
		{
			name: "mm to inch",
			in:   units.Measurement{Quantity: 125, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 4.92, Unit: units.Inch},
		},
		{
			name: "mm to feet",
			in:   units.Measurement{Quantity: 125, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 0.41, Unit: units.Feet},
		},
		{
			name: "mm to yard",
			in:   units.Measurement{Quantity: 125, Unit: units.Millimeter},
			want: units.Measurement{Quantity: 0.14, Unit: units.Yard},
		},
		{
			name: "cm to mm",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 1250, Unit: units.Millimeter},
		},
		{
			name: "cm to cm",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 125, Unit: units.Centimeter},
		},
		{
			name: "cm to m",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 1.25, Unit: units.Meter},
		},
		{
			name: "cm to inch",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 49.21, Unit: units.Inch},
		},
		{
			name: "cm to feet",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 4.1, Unit: units.Feet},
		},
		{
			name: "cm to yard",
			in:   units.Measurement{Quantity: 125, Unit: units.Centimeter},
			want: units.Measurement{Quantity: 1.37, Unit: units.Yard},
		},
		{
			name: "m to mm",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 1250, Unit: units.Millimeter},
		},
		{
			name: "m to cm",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 125, Unit: units.Centimeter},
		},
		{
			name: "m to m",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 1.25, Unit: units.Meter},
		},
		{
			name: "m to inch",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 49.21, Unit: units.Inch},
		},
		{
			name: "m to feet",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 4.1, Unit: units.Feet},
		},
		{
			name: "m to yard",
			in:   units.Measurement{Quantity: 1.25, Unit: units.Meter},
			want: units.Measurement{Quantity: 1.37, Unit: units.Yard},
		},
		{
			name: "inch to mm",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 711.2, Unit: units.Millimeter},
		},
		{
			name: "inch to cm",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 71.12, Unit: units.Centimeter},
		},
		{
			name: "inch to m",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 0.71, Unit: units.Meter},
		},
		{
			name: "inch to inch",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 28, Unit: units.Inch},
		},
		{
			name: "inch to feet",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 2.33, Unit: units.Feet},
		},
		{
			name: "inch to yard",
			in:   units.Measurement{Quantity: 28, Unit: units.Inch},
			want: units.Measurement{Quantity: 0.78, Unit: units.Yard},
		},
		{
			name: "feet to mm",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 853440, Unit: units.Millimeter},
		},
		{
			name: "feet to cm",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 853.44, Unit: units.Centimeter},
		},
		{
			name: "feet to m",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 8.53, Unit: units.Meter},
		},
		{
			name: "feet to inch",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 2.33, Unit: units.Inch},
		},
		{
			name: "feet to feet",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 28, Unit: units.Feet},
		},
		{
			name: "feet to yard",
			in:   units.Measurement{Quantity: 28, Unit: units.Feet},
			want: units.Measurement{Quantity: 9.33, Unit: units.Yard},
		},
		{
			name: "yard to mm",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 3200.4, Unit: units.Millimeter},
		},
		{
			name: "yard to cm",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 320.04, Unit: units.Centimeter},
		},
		{
			name: "yard to m",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 3.2, Unit: units.Meter},
		},
		{
			name: "yard to inch",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 126, Unit: units.Inch},
		},
		{
			name: "yard to feet",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 10.5, Unit: units.Feet},
		},
		{
			name: "yard to yard",
			in:   units.Measurement{Quantity: 3.5, Unit: units.Yard},
			want: units.Measurement{Quantity: 3.5, Unit: units.Yard},
		},
		{
			name: "mg to mg",
			in:   units.Measurement{Quantity: 877, Unit: units.Milligram},
			want: units.Measurement{Quantity: 877, Unit: units.Milligram},
		},
		{
			name: "mg to g",
			in:   units.Measurement{Quantity: 877, Unit: units.Milligram},
			want: units.Measurement{Quantity: 0.88, Unit: units.Gram},
		},
		{
			name: "mg to kg",
			in:   units.Measurement{Quantity: 877, Unit: units.Milligram},
			want: units.Measurement{Quantity: 0.00088, Unit: units.Kilogram},
		},
		{
			name: "mg to lb",
			in:   units.Measurement{Quantity: 877, Unit: units.Milligram},
			want: units.Measurement{Quantity: 0.0019, Unit: units.Pound},
		},
		{
			name: "mg to oz",
			in:   units.Measurement{Quantity: 877, Unit: units.Milligram},
			want: units.Measurement{Quantity: 0.031, Unit: units.Ounce},
		},
		{
			name: "g to mg",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Gram},
			want: units.Measurement{Quantity: 8200, Unit: units.Milligram},
		},
		{
			name: "g to g",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Gram},
			want: units.Measurement{Quantity: 8.2, Unit: units.Gram},
		},
		{
			name: "g to kg",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Gram},
			want: units.Measurement{Quantity: 0.0082, Unit: units.Kilogram},
		},
		{
			name: "g to lb",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Gram},
			want: units.Measurement{Quantity: 0.018, Unit: units.Pound},
		},
		{
			name: "g to oz",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Gram},
			want: units.Measurement{Quantity: 0.29, Unit: units.Ounce},
		},
		{
			name: "kg to mg",
			in:   units.Measurement{Quantity: 1.32, Unit: units.Kilogram},
			want: units.Measurement{Quantity: 1.32e+6, Unit: units.Milligram},
		},
		{
			name: "kg to g",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Kilogram},
			want: units.Measurement{Quantity: 8200, Unit: units.Gram},
		},
		{
			name: "kg to kg",
			in:   units.Measurement{Quantity: 8.2, Unit: units.Kilogram},
			want: units.Measurement{Quantity: 8.2, Unit: units.Kilogram},
		},
		{
			name: "kg to lb",
			in:   units.Measurement{Quantity: 1.32, Unit: units.Kilogram},
			want: units.Measurement{Quantity: 2.91, Unit: units.Pound},
		},
		{
			name: "kg to oz",
			in:   units.Measurement{Quantity: 1.32, Unit: units.Kilogram},
			want: units.Measurement{Quantity: 46.56, Unit: units.Ounce},
		},
		{
			name: "lb to mg",
			in:   units.Measurement{Quantity: 68.2, Unit: units.Pound},
			want: units.Measurement{Quantity: 30934999.63, Unit: units.Milligram},
		},
		{
			name: "lb to g",
			in:   units.Measurement{Quantity: 68.2, Unit: units.Pound},
			want: units.Measurement{Quantity: 30935, Unit: units.Gram},
		},
		{
			name: "lb to kg",
			in:   units.Measurement{Quantity: 68.2, Unit: units.Pound},
			want: units.Measurement{Quantity: 30.94, Unit: units.Kilogram},
		},
		{
			name: "lb to lb",
			in:   units.Measurement{Quantity: 68.2, Unit: units.Pound},
			want: units.Measurement{Quantity: 68.2, Unit: units.Pound},
		},
		{
			name: "lb to oz",
			in:   units.Measurement{Quantity: 68.2, Unit: units.Pound},
			want: units.Measurement{Quantity: 1091.2, Unit: units.Ounce},
		},
		{
			name: "oz to mg",
			in:   units.Measurement{Quantity: 48, Unit: units.Ounce},
			want: units.Measurement{Quantity: 1.36077696e+6, Unit: units.Milligram},
		},
		{
			name: "oz to g",
			in:   units.Measurement{Quantity: 48, Unit: units.Ounce},
			want: units.Measurement{Quantity: 1.36077696e+3, Unit: units.Gram},
		},
		{
			name: "oz to kg",
			in:   units.Measurement{Quantity: 48, Unit: units.Ounce},
			want: units.Measurement{Quantity: 1.36, Unit: units.Kilogram},
		},
		{
			name: "oz to lb",
			in:   units.Measurement{Quantity: 48, Unit: units.Ounce},
			want: units.Measurement{Quantity: 3, Unit: units.Pound},
		},
		{
			name: "oz to oz",
			in:   units.Measurement{Quantity: 48, Unit: units.Ounce},
			want: units.Measurement{Quantity: 48, Unit: units.Ounce},
		},
		{
			name: "mL to mL",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 476, Unit: units.Millilitre},
		},
		{
			name: "mL to L",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 0.476, Unit: units.Litre},
		},
		{
			name: "mL to dL",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 4.76, Unit: units.Decilitre},
		},
		{
			name: "mL to tsp",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 95.2, Unit: units.Teaspoon},
		},
		{
			name: "mL to tbsp",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 31.73, Unit: units.Tablespoon},
		},
		{
			name: "mL to fl oz",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 16.1, Unit: units.FlOz},
		},
		{
			name: "mL to cup",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 2.01, Unit: units.Cup},
		},
		{
			name: "mL to pt",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 1.01, Unit: units.Pint},
		},
		{
			name: "mL to qt",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 0.5, Unit: units.Quart},
		},
		{
			name: "mL to gal",
			in:   units.Measurement{Quantity: 476, Unit: units.Millilitre},
			want: units.Measurement{Quantity: 0.13, Unit: units.Gallon},
		},
		{
			name: "L to mL",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 2360, Unit: units.Millilitre},
		},
		{
			name: "L to L",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 2.36, Unit: units.Litre},
		},
		{
			name: "L to dL",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 23.6, Unit: units.Decilitre},
		},
		{
			name: "L to tsp",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 478.81, Unit: units.Teaspoon},
		},
		{
			name: "L to tbsp",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 159.6, Unit: units.Tablespoon},
		},
		{
			name: "L to fl oz",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 79.8, Unit: units.FlOz},
		},
		{
			name: "L to cup",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 9.98, Unit: units.Cup},
		},
		{
			name: "L to pt",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 4.99, Unit: units.Pint},
		},
		{
			name: "L to qt",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 2.49, Unit: units.Quart},
		},
		{
			name: "L to gal",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Litre},
			want: units.Measurement{Quantity: 0.62, Unit: units.Gallon},
		},
		{
			name: "dL to mL",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 236, Unit: units.Millilitre},
		},
		{
			name: "dL to L",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 0.24, Unit: units.Litre},
		},
		{
			name: "dL to dL",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
		},
		{
			name: "dL to tsp",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 47.88, Unit: units.Teaspoon},
		},
		{
			name: "dL to tbsp",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 15.96, Unit: units.Tablespoon},
		},
		{
			name: "dL to fl oz",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 7.98, Unit: units.FlOz},
		},
		{
			name: "dL to cup",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 0.998, Unit: units.Cup},
		},
		{
			name: "dL to pt",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 0.5, Unit: units.Pint},
		},
		{
			name: "dL to qt",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 0.249, Unit: units.Quart},
		},
		{
			name: "dL to gal",
			in:   units.Measurement{Quantity: 2.36, Unit: units.Decilitre},
			want: units.Measurement{Quantity: 0.62, Unit: units.Gallon},
		},
		{
			name: "tsp to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 115, Unit: units.Millilitre},
		},
		{
			name: "tsp to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 0.12, Unit: units.Litre},
		},
		{
			name: "tsp to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 1.13, Unit: units.Decilitre},
		},
		{
			name: "tsp to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 23, Unit: units.Teaspoon},
		},
		{
			name: "tsp to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 7.67, Unit: units.Tablespoon},
		},
		{
			name: "tsp to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 3.83, Unit: units.FlOz},
		},
		{
			name: "tsp to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 0.48, Unit: units.Cup},
		},
		{
			name: "tsp to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 0.24, Unit: units.Pint},
		},
		{
			name: "tsp to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 0.12, Unit: units.Quart},
		},
		{
			name: "tsp to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Teaspoon},
			want: units.Measurement{Quantity: 0.03, Unit: units.Gallon},
		},
		{
			name: "tbsp to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 340.1, Unit: units.Millilitre},
		},
		{
			name: "tbsp to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 0.34, Unit: units.Litre},
		},
		{
			name: "tbsp to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 3.4, Unit: units.Decilitre},
		},
		{
			name: "tbsp to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 69, Unit: units.Teaspoon},
		},
		{
			name: "tbsp to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 23, Unit: units.Tablespoon},
		},
		{
			name: "tbsp to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 11.5, Unit: units.FlOz},
		},
		{
			name: "tbsp to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 1.44, Unit: units.Cup},
		},
		{
			name: "tbsp to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 0.72, Unit: units.Pint},
		},
		{
			name: "tbsp to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 0.36, Unit: units.Quart},
		},
		{
			name: "tbsp to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Tablespoon},
			want: units.Measurement{Quantity: 0.09, Unit: units.Gallon},
		},
		{
			name: "fl oz to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 680.19, Unit: units.Millilitre},
		},
		{
			name: "fl oz to L",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 0.68, Unit: units.Litre},
		},
		{
			name: "fl oz to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 6.8, Unit: units.Decilitre},
		},
		{
			name: "fl oz to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 138, Unit: units.Teaspoon},
		},
		{
			name: "fl oz to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 46, Unit: units.Tablespoon},
		},
		{
			name: "fl oz to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 23, Unit: units.FlOz},
		},
		{
			name: "fl oz to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 2.88, Unit: units.Cup},
		},
		{
			name: "fl oz to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 1.44, Unit: units.Pint},
		},
		{
			name: "fl oz to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 0.72, Unit: units.Quart},
		},
		{
			name: "fl oz to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.FlOz},
			want: units.Measurement{Quantity: 0.18, Unit: units.Gallon},
		},
		{
			name: "cup to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 5441.53, Unit: units.Millilitre},
		},
		{
			name: "cup to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 5.44, Unit: units.Litre},
		},
		{
			name: "cup to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 54.42, Unit: units.Decilitre},
		},
		{
			name: "cup to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 1104, Unit: units.Teaspoon},
		},
		{
			name: "cup to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 368, Unit: units.Tablespoon},
		},
		{
			name: "cup to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 184, Unit: units.FlOz},
		},
		{
			name: "cup to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 23, Unit: units.Cup},
		},
		{
			name: "cup to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 11.5, Unit: units.Pint},
		},
		{
			name: "cup to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 5.75, Unit: units.Quart},
		},
		{
			name: "cup to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Cup},
			want: units.Measurement{Quantity: 1.44, Unit: units.Gallon},
		},
		{
			name: "pt to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 10883.06, Unit: units.Millilitre},
		},
		{
			name: "pt to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 10.88, Unit: units.Litre},
		},
		{
			name: "pt to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 108.83, Unit: units.Decilitre},
		},
		{
			name: "pt to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 2208, Unit: units.Teaspoon},
		},
		{
			name: "pt to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 736, Unit: units.Tablespoon},
		},
		{
			name: "pt to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 368, Unit: units.FlOz},
		},
		{
			name: "pt to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 46, Unit: units.Cup},
		},
		{
			name: "pt to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 23, Unit: units.Pint},
		},
		{
			name: "pt to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 11.5, Unit: units.Quart},
		},
		{
			name: "pt to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Pint},
			want: units.Measurement{Quantity: 2.88, Unit: units.Gallon},
		},
		{
			name: "qt to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 21766.12, Unit: units.Millilitre},
		},
		{
			name: "qt to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 21.77, Unit: units.Litre},
		},
		{
			name: "qt to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 217.66, Unit: units.Decilitre},
		},
		{
			name: "qt to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 4416, Unit: units.Teaspoon},
		},
		{
			name: "qt to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 1472, Unit: units.Tablespoon},
		},
		{
			name: "qt to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 736, Unit: units.FlOz},
		},
		{
			name: "qt to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 92, Unit: units.Cup},
		},
		{
			name: "qt to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 46, Unit: units.Pint},
		},
		{
			name: "qt to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 23, Unit: units.Quart},
		},
		{
			name: "qt to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Quart},
			want: units.Measurement{Quantity: 5.75, Unit: units.Gallon},
		},
		{
			name: "gal to mL",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 87064.47, Unit: units.Millilitre},
		},
		{
			name: "gal to L",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 87.06, Unit: units.Litre},
		},
		{
			name: "gal to dL",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 870.64, Unit: units.Decilitre},
		},
		{
			name: "gal to tsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 17664, Unit: units.Teaspoon},
		},
		{
			name: "gal to tbsp",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 5888, Unit: units.Tablespoon},
		},
		{
			name: "gal to fl oz",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 2944, Unit: units.FlOz},
		},
		{
			name: "gal to cup",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 368, Unit: units.Cup},
		},
		{
			name: "gal to pt",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 184, Unit: units.Pint},
		},
		{
			name: "gal to qt",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 92, Unit: units.Quart},
		},
		{
			name: "gal to gal",
			in:   units.Measurement{Quantity: 23, Unit: units.Gallon},
			want: units.Measurement{Quantity: 23, Unit: units.Gallon},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.in.Convert(tc.want.Unit)
			if err != nil {
				t.Fatal(err)
			}

			assertEqual(t, got.Unit, tc.want.Unit)
			if math.Abs(got.Quantity-tc.want.Quantity) > 1e-2 {
				t.Errorf("got %#v but want %#v", got.Quantity, tc.want.Quantity)
			}
		})
	}
}

func TestReplaceDecimalFractions(t *testing.T) {
	type testcase struct {
		name string
		in   string
		want string
	}

	decimals := map[string]string{
		"0.50":    "1/2",
		"0.333":   "1/3",
		"0.666":   "2/3",
		"0.25":    "1/4",
		"0.75":    "3/4",
		"0.2":     "1/5",
		"0.4":     "2/5",
		"0.6":     "3/5",
		"0.8":     "4/5",
		"0.16":    "1/6",
		"0.83":    "5/6",
		"0.14":    "1/7",
		"0.125":   "1/8",
		"0.375":   "3/8",
		"0.625":   "5/8",
		"0.875":   "7/8",
		"0.11":    "1/9",
		"0.1":     "1/10",
		"3.33333": "3 1/3",
		"6.5":     "6 1/2",
	}

	var testcases []testcase
	for k, v := range decimals {
		testcases = append(testcases, testcase{
			name: k,
			in:   k + " pineapples " + k,
			want: v + " pineapples " + v,
		})
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := units.ReplaceDecimalFractions(tc.in)
			assertEqual(t, got, tc.want)
		})
	}
}

func TestReplaceVulgarFractions(t *testing.T) {
	type testcase struct {
		name string
		in   string
		want string
	}

	vulgar := map[string]string{
		"½": "1/2",
		"⅓": "1/3",
		"⅔": "2/3",
		"¼": "1/4",
		"¾": "3/4",
		"⅕": "1/5",
		"⅖": "2/5",
		"⅗": "3/5",
		"⅘": "4/5",
		"⅙": "1/6",
		"⅚": "5/6",
		"⅐": "1/7",
		"⅛": "1/8",
		"⅜": "3/8",
		"⅝": "5/8",
		"⅞": "7/8",
		"⅑": "1/9",
		"⅒": "1/10",
	}

	var testcases []testcase
	for k, v := range vulgar {
		testcases = append(testcases, testcase{
			name: k,
			in:   "1" + k + "pineapples" + k,
			want: "1 " + v + " pineapples " + v,
		})
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := units.ReplaceVulgarFractions(tc.in)
			assertEqual(t, got, tc.want)
		})
	}
}

func assertEqual[T string | units.System | units.Unit](t *testing.T, got, want T) {
	if got != want {
		t.Fatalf("got %q but want %q", got, want)
	}
}
