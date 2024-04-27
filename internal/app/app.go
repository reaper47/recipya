package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	BackupPath string // BackupPath is the directory where the data is backed up.
	DBBasePath string // DBBasePath is the directory where the database files are stored.
	ImagesDir  string // ImagesDir is the directory where user images are stored.
	LogsDir    string // LogsDir is the directory where the logs are stored.

	Info = GeneralInfo{
		Version: semver.Version{
			Major: 1,
			Minor: 2,
			Patch: 0,
		},
	} // Info stores general application information.

	FdcDB     = "fdc.db"     // FdcDB is the name of the FDC database.
	RecipyaDB = "recipya.db" // RecipyaDB is the name of Recipya's main database.

	ErrNoUpdate = errors.New("already latest version") // ErrNoUpdate is the error for when the application is up-to-date.
)

// GeneralInfo holds information on the application.
type GeneralInfo struct {
	IsUpdateAvailable   bool
	LastCheckedUpdateAt time.Time
	LastUpdatedAt       time.Time
	Version             semver.Version
}

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
		baseURL := c.Server.URL
		if runtime.GOOS == "windows" {
			baseURL = "http://127.0.0.1"
		}
		return baseURL + ":" + strconv.Itoa(c.Server.Port)
	}

	localAddr := udpAddr()
	if localAddr == nil {
		return c.Server.URL
	}

	if c.Server.Port == 0 {
		c.Server.Port = localAddr.Port
	}
	port := ":" + strconv.Itoa(c.Server.Port)

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

// Update updates the application's configuration.
func (c *ConfigFile) Update(updated ConfigFile) error {
	if c.Server.IsDemo {
		return errors.New("demo disabled")
	}

	c.Email.From = updated.Email.From
	c.Email.SendGridAPIKey = updated.Email.SendGridAPIKey
	c.Integrations.AzureDI.Endpoint = updated.Integrations.AzureDI.Endpoint
	c.Integrations.AzureDI.Key = updated.Integrations.AzureDI.Key
	c.Server.IsAutologin = updated.Server.IsAutologin
	c.Server.IsNoSignups = updated.Server.IsNoSignups
	c.Server.IsProduction = updated.Server.IsProduction

	if os.Getenv("RECIPYA_IS_TEST") == "true" {
		return nil
	}

	if isRunningInDocker() {
		var (
			autologin string
			noSignups string
			isProd    string
		)

		if c.Server.IsAutologin {
			autologin = "true"
		}

		if c.Server.IsNoSignups {
			noSignups = "true"
		}

		if c.Server.IsProduction {
			isProd = "true"
		}

		_ = os.Setenv("RECIPYA_DI_ENDPOINT", c.Integrations.AzureDI.Endpoint)
		_ = os.Setenv("RECIPYA_DI_KEY", c.Integrations.AzureDI.Key)
		_ = os.Setenv("RECIPYA_EMAIL", c.Email.From)
		_ = os.Setenv("RECIPYA_EMAIL_SENDGRID", c.Email.SendGridAPIKey)
		_ = os.Setenv("RECIPYA_SERVER_AUTOLOGIN", autologin)
		_ = os.Setenv("RECIPYA_SERVER_NO_SIGNUPS", noSignups)
		_ = os.Setenv("RECIPYA_SERVER_IS_PROD", isProd)

		return nil
	}

	xb, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(filepath.Dir(DBBasePath), "config.json"), xb, os.ModePerm)
}

func udpAddr() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

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
	AzureDI AzureDI `json:"azureDocumentIntelligence"`
}

// AzureDI holds configuration data for the Azure AI Document Intelligence integration.
type AzureDI struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
}

// PrepareRequest prepares the HTTP request to analyze a document.
func (a AzureDI) PrepareRequest(file io.Reader) (*http.Request, error) {
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	body := map[string]string{
		"base64Source": base64.StdEncoding.EncodeToString(all),
	}

	xb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, a.Endpoint+"/documentintelligence/documentModels/prebuilt-layout:analyze", bytes.NewBuffer(xb))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("_overload", "analyzeDocument")
	q.Add("api-version", "2024-02-29-preview")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Ocp-Apim-Subscription-Key", a.Key)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ConfigServer holds configuration data for the server.
type ConfigServer struct {
	IsAutologin  bool   `json:"autologin"`
	IsDemo       bool   `json:"isDemo"`
	IsNoSignups  bool   `json:"noSignups"`
	IsProduction bool   `json:"isProduction"`
	Port         int    `json:"port"`
	URL          string `json:"url"`
}

// Init initializes the app. This function must be called when the app starts.
// Its name is not *init* so that the function is not executed during the tests.
func Init() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	baseDir := filepath.Join(dir, "Recipya")

	BackupPath = filepath.Join(baseDir, "Backup")
	DBBasePath = filepath.Join(baseDir, "Database")
	ImagesDir = filepath.Join(baseDir, "Images")
	LogsDir = filepath.Join(baseDir, "Logs")

	xs := []string{BackupPath, DBBasePath, ImagesDir, LogsDir}
	for _, s := range xs {
		err = os.MkdirAll(s, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	setup()

	f, err := os.Open(filepath.Join(baseDir, "config.json"))
	if err != nil {
		NewConfig(nil)
	} else {
		NewConfig(f)
	}
	defer f.Close()

	fmt.Printf("File locations:\n\t- Backups: %s\n\t- Database: %s\n\t- Images: %s\n\t- Logs: %s\n", BackupPath, DBBasePath, ImagesDir, LogsDir)
}

// NewConfig initializes the global Config. It can either be populated from environment variables or the configuration file.
func NewConfig(r io.Reader) {
	if r == nil {
		port, _ := strconv.ParseInt(os.Getenv("RECIPYA_SERVER_PORT"), 10, 32)

		// TODO: Remove in v1.3.0
		if os.Getenv("RECIPYA_VISION_ENDPOINT") != "" {
			fmt.Println("The 'RECIPYA_VISION_ENDPOINT' is deprecated. Please use 'RECIPYA_DI_ENDPOINT'.")
		}
		if os.Getenv("RECIPYA_VISION_KEY") != "" {
			fmt.Println("The 'RECIPYA_VISION_KEY' is deprecated. Please use 'RECIPYA_DI_KEY'.")
		}

		Config = ConfigFile{
			Email: ConfigEmail{
				From:           os.Getenv("RECIPYA_EMAIL"),
				SendGridAPIKey: os.Getenv("RECIPYA_EMAIL_SENDGRID"),
			},
			Integrations: ConfigIntegrations{
				AzureDI: AzureDI{
					Endpoint: strings.TrimSuffix(os.Getenv("RECIPYA_DI_ENDPOINT"), "/"),
					Key:      os.Getenv("RECIPYA_DI_KEY"),
				},
			},
			Server: ConfigServer{
				IsAutologin:  os.Getenv("RECIPYA_SERVER_AUTOLOGIN") == "true",
				IsDemo:       os.Getenv("RECIPYA_SERVER_IS_DEMO") == "true",
				IsNoSignups:  os.Getenv("RECIPYA_SERVER_NO_SIGNUPS") == "true",
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
