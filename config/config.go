package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	RecipesDb  string
	RecipesDir string
}

var (
	defaults = map[string]interface{}{
		"recipesDb":  "./dist/recipes.db",
		"recipesDir": "./recipes",
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
