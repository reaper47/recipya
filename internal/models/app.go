package models

import (
	selfupdater "github.com/mJehanno/ghr-self-updater"
	"time"
)

// AppInfo holds information on the application.
type AppInfo struct {
	IsUpdateAvailable   bool
	LastUpdatedAt       time.Time
	LastCheckedUpdateAt time.Time
	Updater             *selfupdater.Updater
}
