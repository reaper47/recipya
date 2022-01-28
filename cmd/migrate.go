package cmd

import (
	"log"
	"os"

	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/driver"
	"github.com/reaper47/recipya/internal/migration"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Short: "upgrade or downgrade the database",
	Long: `Upgrade or downgrade the database.

Upgrade to the next version:
$ ./recipya migrate up

Upgrade to the latest version:
$ ./recipya migrate up -a

Downgrade to the previous version:
$ ./recipya migrate down
`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"up", "down"},
	Run:       migrate,
}

func init() {
	migrateCmd.Flags().BoolP("all", "a", false, "apply all migrations")
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

	isAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		log.Fatalln(err)
	}

	switch args[0] {
	case "up":
		migration.Up(db, isAll)
	case "down":
		migration.Down(db, isAll)
	}
}
