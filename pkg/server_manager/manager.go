package server_manager

import (
	"github.com/docker/docker/client"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
)

type ServerManager struct {
	types.SharedContainer

	docker *client.Client
}

func (m *ServerManager) Start() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("failed to initialize docker: %v", err)
	}

	log.Printf("initialized docker client")
	m.docker = cli
	m.createNetwork()
}
