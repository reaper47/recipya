package main

import (
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/server"
	"github.com/reaper47/recipya/internal/services"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
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
			{
				Name:    "restart",
				Aliases: []string{"r"},
				Usage:   "restarts the web server",
				Action: func(ctx *cli.Context) error {
					id := os.Getpid()
					fmt.Println(id)
					time.Sleep(5 * time.Second)

					self, err := os.Executable()
					if err != nil {
						return err
					}
					args := os.Args
					env := os.Environ()
					// Windows does not support exec syscall.
					if runtime.GOOS == "windows" {
						cmd := exec.Command(self, args[1:]...)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						cmd.Stdin = os.Stdin
						cmd.Env = env
						err := cmd.Run()
						if err == nil {
							os.Exit(0)
						}
						return err
					}
					return syscall.Exec(self, args, env)
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
