package server_manager

import (
	"github.com/docker/docker/client"
	"github.com/vanillaverse/guardian/pkg/types"
	"os"
)

type ServerManager struct {
	types.GuardianD

	docker *client.Client
}

func (m *ServerManager) Start() {
	cli, err := client.NewEnvClient()
	if err != nil {
		m.Error("failed to initialize docker: %v", err)
		os.Exit(1)
	}

	m.Info("initialized docker client")
	m.docker = cli
	m.createNetwork()
}
