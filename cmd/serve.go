package cmd

import (
	"github.com/reaper47/recipe-hunter/core"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the web server",
	Long: `
---------- Help for 'serve' command ------------
|                                              |
|  Starts the web server.                      |
|                                              |
|  The web server is responsible for serving   |
|  the React application.                      |
|______________________________________________|	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
