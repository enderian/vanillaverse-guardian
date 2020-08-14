package daemon

import (
	"context"
	"github.com/kardianos/service"
	"github.com/vanillaverse/guardian/pkg/server_manager"
	"github.com/vanillaverse/guardian/pkg/types"
	"os"
)

type Daemon struct {
	types.GuardianD

	ctx        context.Context
	srvManager *server_manager.ServerManager
}

func (d *Daemon) Start(s service.Service) error {
	d.ctx = context.Background()

	d.Options.ReadFromFiles()
	d.Logger, _ = s.Logger(nil)
	d.srvManager = &server_manager.ServerManager{
		GuardianD: d.GuardianD,
	}

	_ = d.Logger.Infof("starting guardiand with pid %d", os.Getpid())

	d.initializeRedis()
	d.srvManager.Start()
	d.srvManager.PopulateServers()
	return nil
}

func (d *Daemon) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
