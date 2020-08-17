package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

var baseFolder = "."
var serverFlags = []cli.Flag{
	&cli.StringFlag{
		Name:      "socket",
		Usage:     "--socket [unix socket address]",
		TakesFile: true,
		Value:     baseFolder + "/unix.sock",
	},
	&cli.StringFlag{
		Name:     "name",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "flavor",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "version",
		Required: true,
	},
	&cli.BoolFlag{
		Name:     "persist",
		Required: true,
	},
}

func main() {
	app := cli.App{
		Name:        "guardian",
		Description: "Runs commands on the Guardian daemon",
		Usage:       "guardian [subcommand] [opts]",

		Commands: []*cli.Command{
			{
				Name:   "create",
				Flags:  serverFlags,
				Action: createServer,
			},
		},
	}

	_ = app.Run(os.Args)
}

func createServer(c *cli.Context) error {
	// TODO: Write to socket
	return nil
}
