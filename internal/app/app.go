package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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

// Address assembles the server's web address from its URL and host.
func (c *ConfigFile) Address() string {
	addr := c.URL
	if runtime.GOOS == "windows" && strings.Contains(addr, "0.0.0.0") {
		addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
	}

	if c.Port != 0 {
		addr += ":" + strconv.Itoa(c.Port)
	}
	return addr
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
	dir := filepath.Dir(exe)

	xb, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		fmt.Println("The configuration file must be present.")
		fmt.Println("Did you run ./recipya setup?")
		os.Exit(1)
	}

	_, err = os.Stat(filepath.Join(dir, "fdc.db"))
	if os.IsNotExist(err) {
		fmt.Println("The FDC database must be present.")
		fmt.Println("Did you run ./recipya setup?")
		os.Exit(1)
	}

	err = json.Unmarshal(xb, &Config)
	if err != nil {
		panic(err)
	}

	ImagesDir = filepath.Join(dir, "data", "images")
	_, err = os.Stat(ImagesDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(ImagesDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
