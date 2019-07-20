package main

import (
	"log"
	"os"
	gocreate "psutils/pkg/go-create"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "psutils"
	app.Usage = "Helper for common tasks"
	app.Commands = []cli.Command{
		{
			Name:    "go",
			Aliases: []string{"g"},
			Usage:   "helpers",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "setup new project",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "force, f",
							Usage: "'--force true' cleanup existing directory and create new project",
						},
						cli.StringFlag{
							Name:  "path, p",
							Value: "./new-project",
							Usage: "'--path=./test' path to directory where to create project",
						},
					},
					Action: gocreate.Handler,
				},
			},
		},
	}
	app.Version = "0.1.0"

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
