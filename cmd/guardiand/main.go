package main

import (
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
	"github.com/vanillaverse/guardian/pkg/daemon"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
	"os"
)

var baseFolder = "."
var runFlags = []cli.Flag{
	&cli.StringFlag{
		Name:      "config-file",
		Usage:     "--config-file [config file path]",
		TakesFile: true,
		Value:     baseFolder + "/config.yml",
	},
	&cli.StringFlag{
		Name:      "servers-file",
		Usage:     "--servers-file [servers file path]",
		TakesFile: true,
		Value:     baseFolder + "/servers.yml",
	},
	&cli.StringFlag{
		Name:      "socket",
		Usage:     "--socket [unix socket address]",
		TakesFile: true,
		Value:     baseFolder + "/unix.sock",
	},
}

func main() {
	svcConfig := &service.Config{
		Name:        "guardiand",
		DisplayName: "guardiand",
		Description: "The guardian service for VanillaVerse.",
	}

	dmn := &daemon.Daemon{}
	srv, err := service.New(dmn, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:        "guardiand",
		Description: "Installs and runs the guardian daemon",
		Usage:       "guardiand [subcommand] [opts]",

		Flags: runFlags,
		Commands: []*cli.Command{
			{
				Name:  "install",
				Flags: runFlags,
				Action: func(context *cli.Context) error {
					svcConfig.Arguments = os.Args[2:]
					return srv.Install()
				},
			},
			{
				Name: "uninstall",
				Action: func(context *cli.Context) error {
					return srv.Uninstall()
				},
			},
			{
				Name: "start",
				Action: func(context *cli.Context) error {
					return srv.Start()
				},
			},
			{
				Name: "stop",
				Action: func(context *cli.Context) error {
					return srv.Stop()
				},
			},
			{
				Name: "restart",
				Action: func(context *cli.Context) error {
					return srv.Restart()
				},
			},
		},
		Action: func(context *cli.Context) error {
			dmn.Options = &types.Options{
				ConfigFile:  context.String("config-file"),
				ServersFile: context.String("servers-file"),
				UnixSocket:  context.String("socket"),
			}

			err := srv.Run()
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("error occured: %v", err)
	}
}
