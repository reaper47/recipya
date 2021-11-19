package cmd

import (
	"os"

	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/driver"
	"github.com/reaper47/recipya/internal/migration"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate up|down",
	Short: "Upgrade or downgrade Recipya's database",
	Long: `Upgrade or downgrade Recipya's database.

To upgrade the database to the next version:
$ ./recipya migrate up

To downgrade the database to the previous version:
$ ./recipya migrate down
`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"up", "down"},
	Run:       migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrate(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(1)
	}

	dbOptions := config.NewDBOptions("recipya")
	db := driver.ConnectSqlDB(dbOptions.Dsn())
	defer db.Close()

	switch args[0] {
	case "up":
		migration.Up(db)
	case "down":
		migration.Down(db)
	}
}
