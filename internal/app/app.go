package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

var Config ConfigFile

// ConfigFile holds the contents of config.json.
type ConfigFile struct {
	Email        ConfigEmail `json:"email"`
	IsProduction bool        `json:"isProduction"`
	Port         int         `json:"port"`
	URL          string      `json:"url"`
}

// ConfigEmail holds email configuration variables.
type ConfigEmail struct {
	From           string `json:"from"`
	SendGridAPIKey string `json:"sendGridAPIKey"`
}

// ImagesDir is the directory where user images are stored.
var ImagesDir string

// Init initializes the app. This function must be called when the app starts.
// Its name is not *init* so that the function is not executed during the tests.
func Init() {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	xb, err := os.ReadFile(filepath.Join(filepath.Dir(exe), configFileName))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(xb, &Config); err != nil {
		panic(err)
	}

	ImagesDir = filepath.Join(filepath.Dir(exe), "data", "images")
	if _, err := os.Stat(ImagesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(ImagesDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
