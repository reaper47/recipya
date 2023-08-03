package units_test

import (
	"github.com/reaper47/recipya/internal/units"
	"testing"
)

func TestUnit_String(t *testing.T) {
	testcases := []struct {
		in   units.Unit
		want string
	}{
		{units.Celsius, "°C"},
		{units.Centimeter, "cm"},
		{units.Cup, "cup"},
		{units.Decilitre, "dL"},
		{units.Fahrenheit, "°F"},
		{units.Feet, "feet"},
		{units.FlOz, "fl oz"},
		{units.Gallon, "gallon"},
		{units.Gram, "g"},
		{units.Inch, "inch"},
		{units.Kilogram, "kg"},
		{units.Litre, "L"},
		{units.Meter, "m"},
		{units.Milligram, "mg"},
		{units.Millilitre, "mL"},
		{units.Millimeter, "mm"},
		{units.Ounce, "oz"},
		{units.Pint, "pint"},
		{units.Pound, "lb"},
		{units.Quart, "fl qt"},
		{units.Tablespoon, "tbsp"},
		{units.Teaspoon, "tsp"},
		{units.Yard, "yard"},
	}
	for _, tc := range testcases {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.in.String()

			if got != tc.want {
				t.Errorf("got %q but want %q", got, tc.want)
			}
		})
	}
}
