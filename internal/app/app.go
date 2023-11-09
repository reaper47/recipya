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
	port := ":" + strconv.Itoa(c.Port)

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil || isRunningInDocker() {
		addr := c.URL
		if (runtime.GOOS == "windows" || isRunningInDocker()) && strings.Contains(addr, "0.0.0.0") {
			addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
		}

		if c.Port != 0 {
			addr += port
		}
		return addr
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	addr := strings.SplitAfter(c.URL, "://")[0] + localAddr.IP.String()
	if c.Port != 0 {
		return addr + port
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
		err := os.MkdirAll(ImagesDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
