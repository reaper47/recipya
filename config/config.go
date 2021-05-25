package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	RecipesDb  string
	RecipesDir string

	Host string
	Port int
	Wait int
}

var (
	defaults = map[string]interface{}{
		"recipesDb":  "./dist/recipes.db",
		"recipesDir": "./recipes",
		"host":       "0.0.0.0",
		"port":       3000,
		"wait":       15,
	}
	configName  = "config"
	configPaths = []string{".", "./config/"}
	Config      = &config{}
)

func InitConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}
	viper.SetConfigName(configName)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Reading configuration file failed: %v", err)
	}

	err = viper.Unmarshal(Config)
	if err != nil {
		log.Fatalf("Failed to unmarshall the configuration file: %v", err)
	}
}
