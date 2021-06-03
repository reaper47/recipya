package config

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/reaper47/recipe-hunter/consts"
	"github.com/spf13/viper"
)

// ConfigStruct stores the configuration options
// from the config.yaml file.
type ConfigStruct struct {
	RecipesDb     string
	RecipesDir    string
	IndexInterval string

	Host string
	Port int
	Wait int
}

var (
	defaults = map[string]interface{}{
		"recipesDb":     "./bin/recipes.db",
		"recipesDir":    "./recipes",
		"indexInterval": "1d",
		"host":          "0.0.0.0",
		"port":          3000,
		"wait":          15,
	}
	configName  = "config"
	configPaths = []string{".", "./config/"}
)

// Config is the package variable used to
// retrieve configuration options.
var Config      = &ConfigStruct{}

// InitConfig initializes the configuration object with the
// variables from the configuration file or environment.
func InitConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName(configName)

	viper.SetEnvPrefix("RH")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Reading configuration file failed: %v", err)
	}

	err = viper.Unmarshal(Config)
	if err != nil {
		log.Fatalf("Failed to unmarshall the configuration file: %v", err)
	}

	err = Config.Validate()
	if err != nil {
		log.Fatalln(err)
	}
}

// Validate ensures the configuration values are valid.
func (c *ConfigStruct) Validate() error {
	if c.Wait < 1 {
		return ErrWaitNegative
	}

	match, _ := regexp.MatchString("^[1-9]([0-9]?)+[m,h,d,M,w,y]$", c.IndexInterval)
	if !match {
		return ErrIndexIntervalInvalid
	}

	return nil
}

// IndexIntervalToDuration converts the 'indexInterval' config option
// to a time.Duration value.
//
// It is assumed that the interval is valid because the validity
// of the configuration file is checked on startup.
func (c *ConfigStruct) IndexIntervalToDuration() time.Duration {
	r := []rune(c.IndexInterval)
	unit := r[(len(r) - 1)]
	value, _ := strconv.Atoi(string(r[0:(len(r) - 1)]))
	duration := time.Duration(value)

	switch unit {
	case 'm':
		return time.Minute * duration
	case 'h':
		return time.Hour * duration
	case 'd':
		return consts.HoursPerDay * duration
	case 'w':
		return consts.HoursPerWeek * duration
	case 'M':
		return consts.HoursPerMonth * duration
	case 'y':
		duration *= 12
		return consts.HoursPerMonth * duration
	default:
		break
	}

	return consts.HoursPerDay
}
