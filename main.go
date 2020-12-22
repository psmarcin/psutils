package main

import (
	"fmt"
	"os"
	"psutils/pkg/accounting"
	"psutils/pkg/config"
	"sort"
	"time"

	"github.com/prometheus/common/log"
	"github.com/urfave/cli"
)

func main() {
	c := config.Load()
	app := cli.NewApp()
	app.Name = "psutils"
	app.Usage = "Helper for common tasks"
	app.Commands = []cli.Command{
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "manage your config",
			Subcommands: []cli.Command{
				{
					Name:    "edit",
					Aliases: []string{"e"},
					Usage:   "open config file in your default editor",
					Action:  config.HandleEdit,
				},
			},
		},
		{
			Name:    "accounting",
			Aliases: []string{"a"},
			Usage:   "collection of useful commands for managing company",
			Subcommands: []cli.Command{
				{
					Name:    "rename",
					Aliases: []string{"r"},
					Usage:   "reformat pdf file to standard format",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "date, d",
							Value: time.Now().Format(c.Other.MontDateFormat),
							Usage: fmt.Sprintf("date for invoice in format %s, it will be a part of filename - `DATE`-item-faktura.pdf", c.Other.MontDateFormat),
						},
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "name for invoice e.g. vpn-server, it will be a part of filename - 2019-08-`NAME`-faktura.pdf",
						},
						cli.StringFlag{
							Name:  "types, t",
							Value: "faktura", // it's invoice in polish :)
							Usage: "type for invoice e.g. faktura, it will be a part of filename - 2019-08-serwer-vpn-`TYPES`.pdf",
						},
					},
					Action: accounting.Handler,
				},
				{
					Name:    "move",
					Aliases: []string{"m"},
					Usage:   "move file to monthly directory",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "date, d",
							Value: time.Now().Format(c.Other.MontDateFormat),
							Usage: fmt.Sprintf("date for directory in format %s", c.Other.MontDateFormat),
						},
					},
					Action: accounting.MoveHandler,
				},
			},
		},
	}
	app.Version = "1.0.0"

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
