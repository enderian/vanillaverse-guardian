package daemon

import (
	"context"
	"github.com/kardianos/service"
	"github.com/vanillaverse/guardian/pkg/server_manager"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
	"os"
)

type Daemon struct {
	types.SharedContainer

	ctx        context.Context
	srvManager *server_manager.ServerManager
}

func (d *Daemon) Start(s service.Service) error {
	d.ctx = context.Background()

	log.Printf("starting guardiand with pid %d", os.Getpid())

	d.Options.ReadFromFiles()
	d.Logger, _ = s.Logger(nil)
	d.srvManager = &server_manager.ServerManager{
		SharedContainer: d.SharedContainer,
	}

	d.initializeRedis()
	d.srvManager.Start()
	d.srvManager.PopulateServers()
	return nil
}

func (d *Daemon) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
