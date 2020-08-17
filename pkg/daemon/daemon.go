package daemon

import (
	"context"
	"github.com/kardianos/service"
	"github.com/vanillaverse/guardian/pkg/server_manager"
	"github.com/vanillaverse/guardian/pkg/types"
	"net"
	"os"
)

type Daemon struct {
	types.Guardian

	ctx        context.Context
	recv       net.Listener
	srvManager *server_manager.ServerManager
}

func (d *Daemon) Start(s service.Service) error {
	d.ctx = context.Background()

	d.Options.ReadFromFiles()
	d.Logger, _ = s.Logger(nil)

	d.initializeRedis()
	go d.listenRedis()
	go d.initializeSocket()

	d.srvManager = &server_manager.ServerManager{
		Guardian: d.Guardian,
	}
	d.srvManager.Start()
	d.srvManager.PopulateServers()

	_ = d.Logger.Infof("started guardiand with pid %d", os.Getpid())
	return nil
}

func (d *Daemon) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
