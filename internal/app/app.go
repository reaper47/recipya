package app

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	Version        = "1.0.0"
	configFileName = "config.json"
)

// Config references a global ConfigFile.
var Config ConfigFile

// ConfigFile holds the contents of config.json.
type ConfigFile struct {
	Email        ConfigEmail        `json:"email"`
	Integrations ConfigIntegrations `json:"integrations"`
	Server       ConfigServer       `json:"server"`
}

// Address assembles the server's web address from its URL and host.
func (c *ConfigFile) Address() string {
	isLocalhost := strings.Contains(c.Server.URL, "0.0.0.0") ||
		strings.Contains(c.Server.URL, "localhost") ||
		strings.Contains(c.Server.URL, "127.0.0.1")

	if c.Server.IsProduction && !isLocalhost {
		return c.Server.URL
	}

	port := ":" + strconv.Itoa(c.Server.Port)

	if runtime.GOOS == "windows" && isLocalhost {
		return c.Server.URL + port
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil || isRunningInDocker() {
		addr := c.Server.URL
		if (runtime.GOOS == "windows" || isRunningInDocker()) && strings.Contains(addr, "0.0.0.0") {
			addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
		}

		if c.Server.Port != 0 {
			addr += port
		}
		return addr
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	addr := strings.SplitAfter(c.Server.URL, "://")[0] + localAddr.IP.String()
	if c.Server.Port != 0 {
		return addr + port
	}
	return addr
}

// ConfigEmail holds email configuration variables.
type ConfigEmail struct {
	From           string `json:"from"`
	MaxNumberUsers int    `json:"maxNumberUsers"`
	SendGridAPIKey string `json:"sendGridAPIKey"`
}

// ConfigIntegrations holds configuration data for 3rd-party services.
type ConfigIntegrations struct {
	AzureComputerVision AzureComputerVision `json:"azureComputerVision"`
}

// AzureComputerVision holds configuration data for the Azure Computer Vision API.
type AzureComputerVision struct {
	ResourceKey    string `json:"resourceKey"`
	VisionEndpoint string `json:"visionEndpoint"`
}

// ConfigServer holds configuration data for the server.
type ConfigServer struct {
	IsDemo       bool   `json:"isDemo"`
	IsProduction bool   `json:"isProduction"`
	Port         int    `json:"port"`
	URL          string `json:"url"`
}

// ImagesDir is the directory where user images are stored.
var ImagesDir string

// Init initializes the app. This function must be called when the app starts.
// Its name is not *init* so that the function is not executed during the tests.
func Init() {
	setup()

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exe)

	xb, err := os.ReadFile(filepath.Join(dir, configFileName))
	if err != nil {
		fmt.Println("The configuration file must be present.")
		os.Exit(1)
	}

	err = json.Unmarshal(xb, &Config)
	if err != nil {
		panic(err)
	}

	ImagesDir = filepath.Join(dir, "data", "images")
	_, err = os.Stat(ImagesDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(ImagesDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
