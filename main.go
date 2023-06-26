package main

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app.Init()

	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "starts the web server",
				Action: func(ctx *cli.Context) error {
					srv := server.NewServer(services.NewSQLiteService(), services.NewEmailService(), services.NewFilesService())
					srv.Run()
					return nil
				},
			},
		},
		Usage: "the ultimate recipes manager for you and your family",
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
