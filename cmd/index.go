/*
Copyright © 2021 Marc-Andre Charland <macpoule@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/reaper47/recipe-hunter/api"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
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
		api.Index()
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
