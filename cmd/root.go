package cmd

import (
	"github.com/reaper47/recipe-hunter/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "recipe-hunter",
	Short: "Search for what you can cook with the ingredients in your fridge",
	Long: `
                 Recipe Hunter 
	
Recipe Hunter is an application used to search 
for what you can cook with the ingredients in your fridge.

It features a command line interface and a web application.

The user must have a folder of JSON recipes adhering to the 
recipe schema standard (https://schema.org/Recipe). Every 
recipe under this folder will be added and indexed in a 
SQLite database.
`,
}

// Execute runs the root command of the .
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(config.InitConfig)
}
