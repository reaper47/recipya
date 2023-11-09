package units

import "embed"

//go:embed lang/*.json
var fs embed.FS

const (
	Invalid Unit = iota
	Celsius
	Centimeter
	Cup
	Decilitre
	Fahrenheit
	Feet
	FlOz
	Gallon
	Gram
	Inch
	Kilogram
	Litre
	Meter
	Milligram
	Millilitre
	Millimeter
	Ounce
	Pint
	Pound
	Quart
	Tablespoon
	Teaspoon
	Yard
)

// Unit type is a string alias representing a unit.
type Unit int

// String represents the Unit as a string.
func (u Unit) String() string {
	switch u {
	case Celsius:
		return "°C"
	case Centimeter:
		return "cm"
	case Cup:
		return "cup"
	case Decilitre:
		return "dL"
	case Fahrenheit:
		return "°F"
	case Feet:
		return "feet"
	case FlOz:
		return "fl oz"
	case Gallon:
		return "gallon"
	case Gram:
		return "g"
	case Inch:
		return "inch"
	case Kilogram:
		return "kg"
	case Litre:
		return "L"
	case Meter:
		return "m"
	case Milligram:
		return "mg"
	case Millilitre:
		return "mL"
	case Millimeter:
		return "mm"
	case Ounce:
		return "oz"
	case Pint:
		return "pint"
	case Pound:
		return "lb"
	case Quart:
		return "fl qt"
	case Tablespoon:
		return "tbsp"
	case Teaspoon:
		return "tsp"
	case Yard:
		return "yard"
	default:
		return "invalid"
	}
}
