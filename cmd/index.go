package cmd

import (
	"github.com/reaper47/recipya/core"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Indexes the recipes database",
	Long: `
---------- Help for 'index' command ----------
|                                            |
|  Indexes the recipes SQLite database from  |
|  the recipes in the folder specified in    |
|  the configuration file.                   |
|____________________________________________|
 `,
	Run: func(cmd *cobra.Command, args []string) {
		core.Index()
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
