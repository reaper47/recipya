package config_test

import (
	"testing"
	"time"

	"github.com/reaper47/recipya/config"
)

var testConfig = &config.ConfigStruct{
	RecipesDb:     "./bin/recipes.db",
	RecipesDir:    "./recipes",
	IndexInterval: "1d",
	Host:          "0.0.0.0",
	Port:          3000,
	Wait:          15,
	Python:        "python",
	Scraper:       "./tools/scraper/scraper.py",
}

func TestConfig(t *testing.T) {
	t.Run("Configuration is valid", test_Validate_ConfigIsValid)
	t.Run("Provides the good interval duration", test_IndexIntervalToDuration_Valid)

}

func test_Validate_ConfigIsValid(t *testing.T) {
	invalidWaits := []int{0, -1}
	for _, v := range invalidWaits {
		testConfig.Wait = v

		err := testConfig.Validate()
		if err != config.ErrWaitNegative {
			t.Fatalf("Error '%v' unexpectedly thrown for wait: %v", err, v)
		}
	}
	testConfig.Wait = 15

	invalidIntervals := []string{"3.3m", "3,1d", "d", "1W", "1ww", "4Y", "1", "0m", "0y"}
	for _, v := range invalidIntervals {
		testConfig.IndexInterval = v

		err := testConfig.Validate()
		if err != config.ErrIndexIntervalInvalid {
			t.Fatalf("Error '%v' unexpectedly thrown for interval: %v", err, v)
		}
	}

	validIntervals := []string{"3m", "76h", "101d", "777w", "43y"}
	for _, v := range validIntervals {
		testConfig.IndexInterval = v

		err := testConfig.Validate()
		if err == config.ErrIndexIntervalInvalid {
			t.Fatalf("'%v' unexpectedly thrown for interval: %v", err, v)

		}
	}
	testConfig.IndexInterval = "1d"
}

func test_IndexIntervalToDuration_Valid(t *testing.T) {
	intervals := map[string]time.Duration{
		"44m": time.Minute * 44,
		"3h":  time.Hour * 3,
		"10d": time.Hour * 240,
		"7d":  time.Hour * 168,
		"1w":  time.Hour * 168,
		"2w":  time.Hour * 336,
		"22M": time.Hour * 16060,
		"1y":  time.Hour * 8760,
	}

	for interval, expected := range intervals {
		testConfig.IndexInterval = interval
		actual := testConfig.IndexIntervalToDuration()
		if actual != expected {
			t.Fatalf("Interval: '%v' != Expected: '%v'\n", actual, expected)
		}
	}
}
