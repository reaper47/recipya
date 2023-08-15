package models

import "github.com/reaper47/recipya/internal/units"

// User holds data related to a user.
type User struct {
	ID    int64
	Email string
}

// UserSettings holds the user's settings.
type UserSettings struct {
	ConvertAutomatically bool
	MeasurementSystem    units.System
}
