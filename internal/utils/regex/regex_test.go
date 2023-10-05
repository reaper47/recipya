package regex_test

import (
	"github.com/reaper47/recipya/internal/utils/regex"
	"testing"
)

func TestRegex_Decimal(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		xs := []string{
			"0.3335",
			"123.30",
			"1024.6894578",
			".3234234",
		}
		for _, s := range xs {
			t.Run("regex is valid "+s, func(t *testing.T) {
				if !regex.Decimal.MatchString(s) {
					t.Fatal("got false when want true")
				}
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		xs := []string{
			"033333",
			"0.3333.56",
			"1.",
			"1.43a",
			".com@",
			"a1.4444",
			"norway@rocks",
		}
		for _, s := range xs {
			t.Run("regex is invalid "+s, func(t *testing.T) {
				if regex.Email.MatchString(s) {
					t.Error("got true when want false")
				}
			})
		}
	})
}

func TestRegex_Digit(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		xs := []string{
			"2 apples",
			"apples 4",
			"200 apples and oranges",
			"apple number 48 and orange number 2391",
			"2.3 tests",
		}
		for _, s := range xs {
			t.Run("regex is valid "+s, func(t *testing.T) {
				if !regex.Digit.MatchString(s) {
					t.Fatal("got false when want true")
				}
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		xs := []string{
			"two apples",
			"email",
			"apples four",
		}
		for _, s := range xs {
			t.Run("regex is invalid "+s, func(t *testing.T) {
				if regex.Digit.MatchString(s) {
					t.Error("got true when want false")
				}
			})
		}
	})
}

func TestRegex_Email(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		emails := []string{
			"james@bond.com",
			"hello@hello.ca",
			"slave@ukrainia.ua",
			"norway@rocks.no",
		}
		for _, email := range emails {
			t.Run("regex is valid", func(t *testing.T) {
				if !regex.Email.MatchString(email) {
					t.Fatal("got false when want true")
				}
			})
		}
	})

	t.Run("invalid", func(t *testing.T) {
		emails := []string{
			"xyzGmail.com",
			"@gmail.com",
			"email",
			"a@.com",
			".com@",
			"a@",
			"norway@rocks",
		}
		for _, email := range emails {
			t.Run("regex is invalid "+email, func(t *testing.T) {
				if regex.Email.MatchString(email) {
					t.Error("got true when want false")
				}
			})
		}
	})
}

func TestRegex_Quantity(t *testing.T) {
	testcasesValid := []struct{ quantity string }{
		{quantity: "1ml"},
		{quantity: "1mL"},
		{quantity: "15 ml"},
		{quantity: "16 mL"},
		{quantity: "1l"},
		{quantity: "1L"},
		{quantity: "15 l"},
		{quantity: "16 L"},
		{quantity: "1°c"},
		{quantity: "1°f"},
		{quantity: "15 °c"},
		{quantity: "16 °f"},
		{quantity: "1°C"},
		{quantity: "1°F"},
		{quantity: "15 °C"},
		{quantity: "16 °F"},
	}
	for _, tc := range testcasesValid {
		t.Run("regex is valid "+tc.quantity, func(t *testing.T) {
			if !regex.Quantity.MatchString(tc.quantity) {
				t.Error("got false when want true for")
			}
		})
	}

	testcasesInvalid := []struct{ quantity string }{
		{quantity: "ml"},
		{quantity: "mL"},
		{quantity: "l"},
		{quantity: "L"},
		{quantity: "°c"},
		{quantity: "°f"},
		{quantity: "°C"},
		{quantity: "°F"},
		{quantity: "15 mX"},
		{quantity: "\"15mx\""},
	}
	for _, tc := range testcasesInvalid {
		t.Run("regex is invalid "+tc.quantity, func(t *testing.T) {
			if regex.Quantity.MatchString(tc.quantity) {
				t.Errorf("got true when want false for %q", tc.quantity)
			}
		})
	}
}

func TestRegex_Anchor(t *testing.T) {
	t.Run("anchor is valid", func(t *testing.T) {
		a := `<a slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`

		if !regex.Anchor.MatchString(a) {
			t.Error("got false when want true")
		}
	})

	xa := []string{
		`<aa slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`,
		`<aslot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;">`,
		`<a slot="guide-links-primary" href="https://www.youtube.com/about/press/" style="display: none;"`,
	}
	for _, a := range xa {
		t.Run("anchor is invalid ", func(t *testing.T) {
			if regex.Anchor.MatchString(a) {
				t.Errorf("got true when want true for %q", a)
			}
		})
	}
}

func TestRegex_HourMinutes(t *testing.T) {
	xs := []string{
		"45:23",
		"45:00",
		"120:59",
	}
	for _, s := range xs {
		t.Run("regex is valid "+s, func(t *testing.T) {
			if !regex.HourMinutes.MatchString(s) {
				t.Error("got false but want true")
			}
		})
	}

	xs = []string{
		":00",
		"4500",
		"120:60",
		"120 43",
		"120-43",
		"120:",
		"-1:43",
		"10:-43",
		"10:80",
	}
	for _, s := range xs {
		t.Run("regex is invalid "+s, func(t *testing.T) {
			if regex.HourMinutes.MatchString(s) {
				t.Error("got true when want false")
			}
		})
	}
}

func TestRegex_Units(t *testing.T) {
	testcases := []struct {
		name string
		in   string
	}{
		{name: "celsius", in: "1 °C"},
		{name: "celsius", in: "1 c"},
		{name: "celsius", in: "1 celsius"},
		{name: "celsius", in: "1degrees Celsius"},
		{name: "celsius", in: "1 degree celsius"},
		{name: "celsius", in: "1°c"},
		{name: "celsius", in: "1c"},

		{name: "fahrenheit", in: "1 °F"},
		{name: "fahrenheit", in: "1 F"},
		{name: "fahrenheit", in: "1 fahrenheit"},
		{name: "fahrenheit", in: "1degrees Fahrenheit"},
		{name: "fahrenheit", in: "1 degree fahrenheit"},
		{name: "fahrenheit", in: "1°f"},
		{name: "fahrenheit", in: "1F"},

		{name: "mm", in: "1 mM"},
		{name: "mm", in: "1 MM"},
		{name: "mm", in: "1 mm"},
		{name: "mm", in: "1Mm"},
		{name: "mm", in: "1 millimeter"},
		{name: "mm", in: "1 millimeters"},
		{name: "mm", in: "1 millimetres"},
		{name: "mm", in: "1millimetre"},

		{name: "cm", in: "1 cM"},
		{name: "cm", in: "1 cM"},
		{name: "cm", in: "1 cm"},
		{name: "cm", in: "1Cm"},
		{name: "cm", in: "1 centimeter"},
		{name: "cm", in: "1 centimeters"},
		{name: "cm", in: "1 centimetres"},
		{name: "cm", in: "1centimetre"},

		{name: "m", in: "1 M"},
		{name: "m", in: "1 m"},
		{name: "m", in: "1m"},
		{name: "m", in: "1 meter"},
		{name: "m", in: "1 meters"},
		{name: "m", in: "1 metres"},
		{name: "m", in: "1metre"},

		{name: "inch", in: "1 in"},
		{name: "inch", in: "1in"},
		{name: "inch", in: "1inch"},
		{name: "inch", in: "1 inches"},
		{name: "inch", in: `1"`},
		{name: "inch", in: `1 "`},

		{name: "ft", in: "1 ft"},
		{name: "ft", in: "1ft"},
		{name: "ft", in: "1foot"},
		{name: "ft", in: "1 feet"},
		{name: "ft", in: "1′"},
		{name: "ft", in: `1 ′`},

		{name: "yard", in: "1 yard"},
		{name: "yard", in: "1yards"},

		{name: "mg", in: "1 mg"},
		{name: "mg", in: "1mG"},
		{name: "mg", in: "1 milligram"},
		{name: "mg", in: "1milligrams"},
		{name: "mg", in: "1 milligramme"},
		{name: "mg", in: "1milligrammes"},

		{name: "g", in: "1 g"},
		{name: "g", in: "1G"},
		{name: "g", in: "1 gram"},
		{name: "g", in: "1grams"},
		{name: "g", in: "1 gramme"},
		{name: "g", in: "1grammes"},

		{name: "kg", in: "1 kg"},
		{name: "kg", in: "1kG"},
		{name: "kg", in: "1 kilogram"},
		{name: "kg", in: "1kilograms"},
		{name: "kg", in: "1 kilogramme"},
		{name: "kg", in: "1kilogrammes"},

		{name: "lb", in: "1 lb"},
		{name: "lb", in: "1#"},
		{name: "lb", in: "1 pound"},
		{name: "lb", in: "1pounds"},
		{name: "lb", in: "1lbs"},
		{name: "lb", in: "1lbs."},
		{name: "lb", in: "1lb."},

		{name: "oz", in: "1 oz"},
		{name: "oz", in: "1oz."},
		{name: "oz", in: "1 ounce"},
		{name: "oz", in: "1ounces"},

		{name: "mL", in: "1 ml"},
		{name: "mL", in: "1 mL"},
		{name: "mL", in: "1 milliliter"},
		{name: "mL", in: "1ml"},
		{name: "mL", in: "1mL"},
		{name: "mL", in: "1millilitres"},
		{name: "mL", in: "1 milliliters"},
		{name: "mL", in: "1millilitre"},

		{name: "dL", in: "1 dl"},
		{name: "dL", in: "1 dL"},
		{name: "dL", in: "1 deciliter"},
		{name: "dL", in: "1dl"},
		{name: "dL", in: "1dL"},
		{name: "dL", in: "1decilitres"},
		{name: "dL", in: "1 deciliters"},
		{name: "dL", in: "1decilitre"},

		{name: "L", in: "1 l"},
		{name: "L", in: "1 L"},
		{name: "L", in: "1 liter"},
		{name: "L", in: "1l"},
		{name: "L", in: "1L"},
		{name: "L", in: "1litres"},
		{name: "L", in: "1 liters"},
		{name: "L", in: "1litre"},

		{name: "tsp", in: "1 tsp"},
		{name: "tsp", in: "1 tsp."},
		{name: "tsp", in: "1 teaspoons"},
		{name: "tsp", in: "1tsp"},
		{name: "tsp", in: "1tsp."},
		{name: "tsp", in: "1teaspoon"},
		{name: "tsp", in: "1teaspoons"},

		{name: "tbsp", in: "1 tbsp"},
		{name: "tbsp", in: "1 tbsp."},
		{name: "tbsp", in: "1 tablespoons"},
		{name: "tbsp", in: "1tbsp"},
		{name: "tbsp", in: "1tbsp."},
		{name: "tbsp", in: "1tablespoon"},
		{name: "tbsp", in: "1tablespoons"},

		{name: "fl oz", in: "1 fl oz"},
		{name: "fl oz", in: "1 floz."},
		{name: "fl oz", in: "1 fluid ounces"},
		{name: "fl oz", in: "1fluid oz"},
		{name: "fl oz", in: "1fluid oz."},
		{name: "fl oz", in: "1fl. oz."},

		{name: "cup", in: "1cup"},
		{name: "cup", in: "1 cup."},
		{name: "cup", in: "1 cups"},
		{name: "cup", in: "1cups"},

		{name: "pint", in: "1pint"},
		{name: "pint", in: "1 pint"},
		{name: "pint", in: "1 pints"},
		{name: "pint", in: "1pints"},
		{name: "pint", in: "1pt"},
		{name: "pint", in: "1 fl.pt."},
		{name: "pint", in: "1 pt."},
		{name: "pint", in: "1 fl. pt."},

		{name: "quart", in: "1quart"},
		{name: "quart", in: "1 quart"},
		{name: "quart", in: "1 quarts"},
		{name: "quart", in: "1quarts"},
		{name: "quart", in: "1qt"},
		{name: "quart", in: "1 fl.qt."},
		{name: "quart", in: "1 qt."},
		{name: "quart", in: "1 fl. qt."},

		{name: "gallon", in: "1gallons"},
		{name: "gallon", in: "1gallon"},
		{name: "gallon", in: "1 gals"},
		{name: "gallon", in: "1 gal"},
		{name: "gallon", in: "1gal"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if !regex.Unit.MatchString(tc.in) {
				t.Fail()
			}
		})
	}
}

func TestRegex_UnitImperial(t *testing.T) {
	valid := []struct{ in string }{
		{"2 cups yay"}, {"1 cup yay"},
		{"3 feet yay"}, {"1 foot yay"}, {"2 ft yay"}, {"3 ft. yay"}, {"3ft yay"}, {"3ft. yay"},
		{"1 fluidounce"}, {"2 fluidounces"}, {"1 fluid oz."}, {"1 fluidoz."}, {"1 fl.oz."}, {"1 fl. oz."},
		{"2 gallons"}, {"1 gallon"}, {"3 gals yay"}, {"1 gal yay"},
		{"cut 1 in"}, {"cut 1 inch"}, {`cut 4"`}, {"cut 5 inches"}, {"cut 5 inche"},
		{"5 ounces"}, {"6 ounce"}, {"5 oz"}, {"5 oz."},
		{"1 pint"}, {"2 pints"}, {"3 pt"}, {"4 pt."},
		{"2 pounds"}, {"1 pound"}, {"1 lb"}, {"1lb"}, {"2lbs"}, {"1lb"}, {"1#"}, {"4 #"},
		{"1 quart"}, {"2 quarts"}, {"3 qt"}, {"4 qt."},
		{"5 tablespoon"}, {"5 tablespoons"}, {"4 tbsp."}, {"4 tbsp"}, {"4tbsp."}, {"4tbsp"},
		{"5 teaspoon"}, {"5 teaspoons"}, {"4 tsp."}, {"4 tsp"}, {"4tsp."}, {"4tsp"},
		{"1 yard"}, {"1yard"}, {"5 yards"}, {"5yards"},
		{"275°f"}, {"275 °f"}, {"275 degrees fahrenheit"}, {"275 fahrenheit"}, {"275 degree fahrenheit"},
	}
	for _, tc := range valid {
		t.Run(tc.in, func(t *testing.T) {
			if !regex.UnitImperial.MatchString(tc.in) {
				t.Fatal("got false when should match")
			}
		})
	}

	invalid := []struct{ in string }{
		{"10 centimeters"}, {"1 centimeter"}, {"1cm"}, {"1 cm"},
		{"10 decilitres"}, {"10 deciliters"}, {"1dl"}, {"10 dl"}, {"5 decilitre"}, {"5 decilitre"},
		{"10 millilitres"}, {"10 milliliters"}, {"1 millilitre"}, {"1 milliliter"}, {"5ml"}, {"5 ml"},
		{"10 millimetres"}, {"10 millimeters"}, {"1 millimetre"}, {"1 millimeter"}, {"5mm"}, {"5 mm"},
		{"10 grams"}, {"10 grammes"}, {"1 gramme"}, {"10 gram"}, {"10 g"}, {"10g"},
		{"10 kilograms"}, {"10 kilogrammes"}, {"1 kilogramme"}, {"10 kilogram"}, {"10 kg"}, {"10kg"},
		{"10 milligrams"}, {"10 milligrammes"}, {"1 milligramme"}, {"10 milligram"}, {"10 mg"}, {"10mg"},
		{"10 metres"}, {"10 meters"}, {"1 metre"}, {"1 meter"}, {"5m"}, {"5 m"},
		{"10 litres"}, {"10 liters"}, {"1 litre"}, {"1 liter"}, {"5l"}, {"5 l"},
		{"275°c"}, {"275 °c"}, {"275 degrees celsius"}, {"275 celsius"}, {"275 degree celsius"},
	}
	for _, tc := range invalid {
		t.Run(tc.in, func(t *testing.T) {
			if regex.UnitImperial.MatchString(tc.in) {
				t.Fatal("got a match when it should not have")
			}
		})
	}
}

func TestRegex_UnitMetric(t *testing.T) {
	valid := []struct{ in string }{
		{"10 centimeters"}, {"1 centimeter"}, {"1cm"}, {"1 cm"},
		{"10 decilitres"}, {"10 deciliters"}, {"1dl"}, {"10 dl"}, {"5 decilitre"}, {"5 decilitre"},
		{"10 millilitres"}, {"10 milliliters"}, {"1 millilitre"}, {"1 milliliter"}, {"5ml"}, {"5 ml"},
		{"10 millimetres"}, {"10 millimeters"}, {"1 millimetre"}, {"1 millimeter"}, {"5mm"}, {"5 mm"},
		{"10 grams"}, {"10 grammes"}, {"1 gramme"}, {"10 gram"}, {"10 g"}, {"10g"},
		{"10 kilograms"}, {"10 kilogrammes"}, {"1 kilogramme"}, {"10 kilogram"}, {"10 kg"}, {"10kg"},
		{"10 milligrams"}, {"10 milligrammes"}, {"1 milligramme"}, {"10 milligram"}, {"10 mg"}, {"10mg"},
		{"10 metres"}, {"10 meters"}, {"1 metre"}, {"1 meter"}, {"5m"}, {"5 m"},
		{"10 litres"}, {"10 liters"}, {"1 litre"}, {"1 liter"}, {"5l"}, {"5 l"},
		{"275°c"}, {"275 °c"}, {"275 degrees celsius"}, {"275 celsius"}, {"275 degree celsius"},
	}
	for _, tc := range valid {
		t.Run(tc.in, func(t *testing.T) {
			if !regex.UnitMetric.MatchString(tc.in) {
				t.Fatal("got false when should match")
			}
		})
	}

	invalid := []struct{ in string }{
		{"2 cups yay"}, {"1 cup yay"},
		{"3 feet yay"}, {"1 foot yay"}, {"2 ft yay"}, {"3 ft. yay"}, {"3ft yay"}, {"3ft. yay"},
		{"1 fluidounce"}, {"2 fluidounces"}, {"1 fluid oz."}, {"1 fluidoz."}, {"1 fl.oz."}, {"1 fl. oz."},
		{"2 gallons"}, {"1 gallon"}, {"3 gals yay"}, {"1 gal yay"},
		{"cut 1 in"}, {"cut 1 inch"}, {`cut 4"`}, {"cut 5 inches"}, {"cut 5 inche"},
		{"5 ounces"}, {"6 ounce"}, {"5 oz"}, {"5 oz."},
		{"1 pint"}, {"2 pints"}, {"3 pt"}, {"4 pt."},
		{"2 pounds"}, {"1 pound"}, {"1 lb"}, {"1lb"}, {"2lbs"}, {"1lb"}, {"1#"}, {"4 #"},
		{"1 quart"}, {"2 quarts"}, {"3 qt"}, {"4 qt."},
		{"5 tablespoon"}, {"5 tablespoons"}, {"4 tbsp."}, {"4 tbsp"}, {"4tbsp."}, {"4tbsp"},
		{"5 teaspoon"}, {"5 teaspoons"}, {"4 tsp."}, {"4 tsp"}, {"4tsp."}, {"4tsp"},
		{"1 yard"}, {"1yard"}, {"5 yards"}, {"5yards"},
		{"275°f"}, {"275 °f"}, {"275 degrees fahrenheit"}, {"275 fahrenheit"}, {"275 degree fahrenheit"},
	}
	for _, tc := range invalid {
		t.Run(tc.in, func(t *testing.T) {
			if regex.UnitMetric.MatchString(tc.in) {
				t.Fatal("got a match when it should not have")
			}
		})
	}
}

func TestRegex_Website(t *testing.T) {
	urls := []string{
		"http://subdomain.example.com",
		"http://www.example.com/path/to/page/subpage",
		"https://www.example.com/path/to/page/subpage?key1=value1&key2=value2",
		"http://username:password@subdomain.example.com:8080/path/to/page/subpage?key1=value1&key2=value2#fragment",
		"http://www.example.com/path/to/page/subpage/subsubpage?key1=value1&key2=value2#fragment",
		"https://subdomain.example.com/path/to/page/subpage/subsubpage?key1=value1&key2=value2#fragment",
		"https://realfood.tesco.com/recipes/salted-honey-and-rosemary-lamb-with-roasties-and-rainbow-carrots.html",
	}

	for _, url := range urls {
		if !regex.URL.MatchString(url) {
			t.Errorf("%v should be matched", url)
		}
	}
}
