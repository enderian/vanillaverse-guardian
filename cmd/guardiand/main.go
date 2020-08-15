package main

import (
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
	"github.com/vanillaverse/guardian/pkg/daemon"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
	"os"
)

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

	runFlags := []cli.Flag{
		&cli.StringFlag{
			Name:      "config-file",
			Usage:     "--config-file [config file path]",
			TakesFile: true,
			Value:     "./config.yml",
		},
		&cli.StringFlag{
			Name:      "servers-file",
			Usage:     "--servers-file [servers file path]",
			TakesFile: true,
			Value:     "./servers.yml",
		},
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
					err := srv.Install()
					if err != nil {
						log.Printf("error while installing guardiand: %v", err)
					}
					return err
				},
			},
			{
				Name: "uninstall",
				Action: func(context *cli.Context) error {
					svcConfig.Arguments = os.Args[2:]
					err := srv.Uninstall()
					if err != nil {
						log.Printf("error while uninstalling guardiand: %v", err)
					}
					return err
				},
			},
			{
				Name: "start",
				Action: func(context *cli.Context) error {
					svcConfig.Arguments = os.Args[2:]
					err := srv.Start()
					if err != nil {
						log.Printf("error while starting guardiand: %v", err)
					}
					return err
				},
			},
			{
				Name: "stop",
				Action: func(context *cli.Context) error {
					svcConfig.Arguments = os.Args[2:]
					err := srv.Stop()
					if err != nil {
						log.Printf("error while stopping guardiand: %v", err)
					}
					return err
				},
			},
			{
				Name: "restart",
				Action: func(context *cli.Context) error {
					svcConfig.Arguments = os.Args[2:]
					err := srv.Restart()
					if err != nil {
						log.Printf("error while restarting guardiand: %v", err)
					}
					return err
				},
			},
		},

		Action: func(context *cli.Context) error {
			dmn.Options = &types.Options{
				ConfigFile:  context.String("config-file"),
				ServersFile: context.String("servers-file"),
			}

			err := srv.Run()
			if err != nil {
				log.Printf("error while running guardiand: %v", err)
			}
			return err
		},
	}

	_ = app.Run(os.Args)
}
