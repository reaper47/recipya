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
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search ingredient1,ingredient2,...,ingredientN",
	Short: "Search for recipes based on ingredients",
	Long: `
---------- Help for 'search' command ---------
|                                            |
|  Searches for recipes based on the list    |
|  of ingredients provided as input.         |
|                                            |
|  In other words, the command returns       |
|  recipes you can prepare or cook from      |
|  ingredients in your fridge.               |
|____________________________________________|	
	 `,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("search called with ingredients: %v\n", args)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
