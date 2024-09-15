package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/reaper47/recipya/internal/app"
)

func TestConfigFile_Address(t *testing.T) {
	testcases := []struct {
		name           string
		in             app.ConfigFile
		want           string
		wantRandomPort bool
	}{
		{
			name:           "without port",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "https://localhost"}},
			want:           "https://localhost",
			wantRandomPort: true,
		},
		{
			name:           "with port",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "http://127.0.0.1", Port: 8078}},
			want:           "http://127.0.0.1:8078",
			wantRandomPort: false,
		},
		{
			name:           "hosted somewhere",
			in:             app.ConfigFile{Server: app.ConfigServer{URL: "https://recipya.com", Port: 8078, IsProduction: true}},
			want:           "https://recipya.com",
			wantRandomPort: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Address()
			if got != tc.want {
				if tc.wantRandomPort {
					split := strings.Split(got, ":")
					port, err := strconv.ParseInt(split[len(split)-1], 10, 64)
					if err != nil || port == 0 {
						t.Fatal("port should not be 0")
						return
					}
				} else {
					t.Fatalf("got %q but want %q", got, tc.want)
				}
			}
			if tc.wantRandomPort && tc.in.Server.Port == 0 {
				t.Fatalf("expected random port %t %d", tc.wantRandomPort, app.Config.Server.Port)
			}
		})
	}
}

func TestConfigServer_IsCookieSecure(t *testing.T) {
	testcases := []struct {
		name     string
		in       app.ConfigFile
		isSecure bool
	}{
		{
			name:     "is prod localhost",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "http://localhost"}},
			isSecure: true,
		},
		{
			name:     "is demo and prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: true, IsProduction: true, Port: 8078, URL: "http://localhost"}},
			isSecure: true,
		},
		{
			name:     "is prod not localhost",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "http://192.168.0.1"}},
			isSecure: false,
		},
		{
			name:     "is not prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: false, Port: 8078, URL: "http://192.168.0.1"}},
			isSecure: false,
		},
		{
			name:     "is demo not prod",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: true, IsProduction: false, Port: 8078, URL: "http://localhost"}},
			isSecure: false,
		},
		{
			name:     "is hosted website",
			in:       app.ConfigFile{Server: app.ConfigServer{IsDemo: false, IsProduction: true, Port: 8078, URL: "https://www.recipya.com"}},
			isSecure: true,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.in.IsCookieSecure() != tc.isSecure {
				if tc.isSecure {
					t.Fatal("should have been secure")
				} else {
					t.Fatal("should not have been secure")
				}
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	base := app.ConfigFile{
		Email: app.ConfigEmail{
			From:           "my@email.com",
			SendGridAPIKey: "API_KEY",
		},
		Integrations: app.ConfigIntegrations{
			AzureDI: app.AzureDI{
				Endpoint: "https://{resource_di}.cognitiveservices.azure.com",
				Key:      "KEY_1",
			},
		},
		Server: app.ConfigServer{
			IsDemo:       false,
			IsProduction: false,
			Port:         8078,
			URL:          "http://0.0.0.0",
		},
	}

	env := map[string]string{
		"RECIPYA_DI_ENDPOINT":    "https://{resource_di}.cognitiveservices.azure.com",
		"RECIPYA_DI_KEY":         "KEY_1",
		"RECIPYA_EMAIL":          "my@email.com",
		"RECIPYA_EMAIL_SENDGRID": "API_KEY",
		"RECIPYA_SERVER_IS_DEMO": "false",
		"RECIPYA_SERVER_IS_PROD": "false",
		"RECIPYA_SERVER_PORT":    "8078",
	}

	t.Run("load from config file", func(t *testing.T) {
		defer func() {
			app.Config = app.ConfigFile{}
		}()
		xb, _ := json.Marshal(&base)

		app.NewConfig(bytes.NewBuffer(xb))
		got := app.Config

		if !cmp.Equal(got, base) {
			t.Log(cmp.Diff(got, base))
			t.Fail()
		}
	})

	t.Run("load from env", func(t *testing.T) {
		defer func() {
			app.Config = app.ConfigFile{}
			os.Clearenv()
		}()
		for k, v := range env {
			_ = os.Setenv(k, v)
		}

		app.NewConfig(nil)
		got := app.Config

		if !cmp.Equal(got, base) {
			t.Log(cmp.Diff(got, base))
			t.Fail()
		}
	})
}

func TestAzureDI_PrepareRequest(t *testing.T) {
	c := app.AzureDI{
		Endpoint: "https://di-rocks.cognitiveservices.azure.com",
		Key:      "trump-farts-in-court",
	}
	body := "Hello"

	got, _ := c.PrepareRequest(strings.NewReader(body))

	if got.Method != http.MethodPost {
		t.Fatalf("got method %q but want POST", got.Method)
	}
	if got.Host != "di-rocks.cognitiveservices.azure.com" {
		t.Fatalf("got url.host %q but want https://di-rocks.cognitiveservices.azure.com", got.Host)
	}
	if got.URL.Path != "/documentintelligence/documentModels/prebuilt-layout:analyze" {
		t.Fatalf("got url.path %q but want /documentintelligence/documentModels/prebuilt-layout:analyze", got.URL.Path)
	}
	if got.URL.RawQuery != "_overload=analyzeDocument&api-version=2024-02-29-preview" {
		t.Fatalf("got url.rawquery %q but want _overload=analyzeDocument&api-version=2024-02-29-preview", got.URL.RawQuery)
	}
	b, _ := io.ReadAll(got.Body)
	if string(b) != `{"base64Source":"SGVsbG8="}` {
		t.Fatalf("got content %q, but want %q", b, body)
	}
	headers := map[string]string{
		"Ocp-Apim-Subscription-Key": "trump-farts-in-court",
		"Content-Type":              "application/json",
	}
	for k, v := range headers {
		actual := got.Header.Get(k)
		if actual != v {
			t.Errorf("got header %q but want header %q", actual, k)
		}
	}
}
