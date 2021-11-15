package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/reaper47/recipya/core"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search ingredient1,ingredient2,...,ingredientN",
	Example: "recipya search avocado,garlic -m 1 -n 5",
	Short:   "Search for recipes based on ingredients",
	Long: `
---------- Help for 'search' command -----------
|                                              |
|  Searches for recipes based on the list of   |
|  ingredients provided as input.              |
|                                              |
|  In other words, the command returns recipes |
|  you can prepare or cook from ingredients in |
|  your fridge.                                |
|                                              |
|  Search modes:                               |
|   1 -> Minimize the number of missing        |
|        ingredients in order to buy less      |
|        at the grocery store.                 |
|   2 -> Maximize the number of ingredients    |
|        taken from the fridge [default]       |
|______________________________________________|
	 `,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ingredients := strings.Split(args[0], ",")
		limit, _ := cmd.Flags().GetInt("num-recipes")
		if limit > 30 {
			limit = 30
		}
		mode, _ := cmd.Flags().GetInt("mode")

		recipes, err := core.Search(ingredients, mode, limit)
		if err != nil {
			log.Fatalln(err)
		}

		for _, r := range recipes {
			fmt.Printf("%+v\n\n", r)
		}
	},
}

func init() {
	searchCmd.Flags().IntP("num-recipes", "n", 10, "number of recipes to return (max: 30)")
	searchCmd.Flags().IntP("mode", "m", 2, "the search mode to employ (see help for options)")
	rootCmd.AddCommand(searchCmd)
}
