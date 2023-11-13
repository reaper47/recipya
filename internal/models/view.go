package models

// ViewMode is an integer type alias for various data view modes.
type ViewMode int

// These constants enumerate all possible viewing modes.
const (
	GridViewMode ViewMode = iota
	ListViewMode
)

// ViewModeFromInt returns the ViewMode for the respective integer.
func ViewModeFromInt(num int64) ViewMode {
	switch num {
	case 0:
		return GridViewMode
	case 1:
		return ListViewMode
	default:
		return GridViewMode
	}
}

// ViewModeFromString returns the ViewMode for the respective string.
func ViewModeFromString(s string) ViewMode {
	switch s {
	case "grid":
		return GridViewMode
	case "list":
		return ListViewMode
	default:
		return GridViewMode
	}
}
