package units

import "strings"

const (
	ImperialSystem System = "imperial"
	InvalidSystem  System = "invalid"
	MetricSystem   System = "metric"
)

// System is a type alias for a measurement system.
type System string

// String represents the MeasurementSystem as a string.
func (m System) String() string {
	return string(m)
}

// NewSystem converts the system string to a MeasurementSystem.
func NewSystem(system string) System {
	switch strings.ToLower(system) {
	case "imperial":
		return ImperialSystem
	case "metric":
		return MetricSystem
	default:
		return InvalidSystem
	}
}
