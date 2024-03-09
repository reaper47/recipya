package models

import "time"

// AppInfo holds information on the application.
type AppInfo struct {
	IsUpdateAvailable   bool
	LastUpdatedAt       time.Time
	LastCheckedUpdateAt time.Time
}
