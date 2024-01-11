package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	Version        = "1.1.0" // Version represents the current version of the application.
	configFileName = "config.json"
)

// ImagesDir is the directory where user images are stored.
var ImagesDir string

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

	if isLocalhost && c.Server.Port > 0 {
		return c.Server.URL + ":" + strconv.Itoa(c.Server.Port)
	}

	localAddr := udpAddr()
	if localAddr == nil {
		return c.Server.URL
	}

	if c.Server.Port == 0 {
		c.Server.Port = localAddr.Port
	}
	port := ":" + strconv.Itoa(c.Server.Port)

	if runtime.GOOS == "windows" && isLocalhost {
		return c.Server.URL + port
	}

	if isRunningInDocker() {
		addr := c.Server.URL
		if runtime.GOOS == "windows" && strings.Contains(addr, "0.0.0.0") {
			addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
		}
		return addr + port
	}

	xs := strings.SplitAfter(c.Server.URL, "://")
	protocol := "https://"
	if len(xs) > 0 {
		protocol = xs[0]
	}

	addr := protocol + localAddr.IP.String()
	if c.Server.Port != 0 {
		return addr + port
	}
	return addr
}

func udpAddr() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer func() {
		_ = conn.Close()
	}()

	return conn.LocalAddr().(*net.UDPAddr)
}

// IsCookieSecure returns whether the cookie should secure.
func (c *ConfigFile) IsCookieSecure() bool {
	u, err := url.ParseRequestURI(c.Server.URL)
	if u == nil || err != nil {
		return false
	}

	host := u.Hostname()
	return c.Server.IsProduction && (u.Scheme == "https" || (host == "localhost" || host == "127.0.0.1"))
}

// ConfigEmail holds email configuration variables.
type ConfigEmail struct {
	From           string `json:"from"`
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

// Init initializes the app. This function must be called when the app starts.
// Its name is not *init* so that the function is not executed during the tests.
func Init() {
	setup()

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exe)

	f, _ := os.Open(filepath.Join(dir, configFileName))
	NewConfig(f)
	_ = f.Close()

	fmt.Printf("%+v\n", Config)

	ImagesDir = filepath.Join(dir, "data", "images")
	_, err = os.Stat(ImagesDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(ImagesDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

// NewConfig initializes the global Config. It can either be populated from environment variables or the configuration file.
func NewConfig(r io.Reader) {
	if r == nil {
		port, _ := strconv.ParseInt(os.Getenv("RECIPYA_SERVER_PORT"), 10, 64)
		Config = ConfigFile{
			Email: ConfigEmail{
				From:           os.Getenv("RECIPYA_EMAIL"),
				SendGridAPIKey: os.Getenv("RECIPYA_EMAIL_SENDGRID"),
			},
			Integrations: ConfigIntegrations{
				AzureComputerVision: AzureComputerVision{
					ResourceKey:    os.Getenv("RECIPYA_VISION_KEY"),
					VisionEndpoint: os.Getenv("RECIPYA_VISION_ENDPOINT"),
				},
			},
			Server: ConfigServer{
				IsDemo:       os.Getenv("RECIPYA_SERVER_IS_DEMO") == "true",
				IsProduction: os.Getenv("RECIPYA_SERVER_IS_PROD") == "true",
				Port:         int(port),
				URL:          os.Getenv("RECIPYA_SERVER_URL"),
			},
		}
	} else {
		err := json.NewDecoder(r).Decode(&Config)
		if err != nil {
			fmt.Println("The configuration file must be present.")
			os.Exit(1)
		}
	}

	if Config.Server.URL == "" {
		fmt.Println("Missing 'server.url' in the configuration.")
		fmt.Println("If you use Docker, please pass the `RECIPYA_SERVER_URL` environment variable.")
		fmt.Println("Otherwise, please double-check your configuration file.")
		os.Exit(1)
	}
}
