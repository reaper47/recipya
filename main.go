// Package main is Recipya's entrypoint.
package main

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/server"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cliApp := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "starts the web server",
				Action: func(_ *cli.Context) error {
					app.Init()
					srv := server.NewServer()
					srv.Run()
					return nil
				},
			},
		},
		Usage: "the ultimate recipes manager for you and your family",
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
