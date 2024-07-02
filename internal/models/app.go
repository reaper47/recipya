package models

import (
	"time"
)

// AppInfo holds information on the application.
type AppInfo struct {
	IsUpdateAvailable   bool
	LastUpdatedAt       time.Time
	LastCheckedUpdateAt time.Time
}

// Replace holds the old and new values for a strings.ReplaceAll operation.
type Replace struct {
	Old string
	New string
}
