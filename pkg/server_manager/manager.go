package server_manager

import (
	"github.com/docker/docker/client"
	"github.com/vanillaverse/guardian/pkg/types"
)

type ServerManager struct {
	types.Guardian

	docker *client.Client
}

func (m *ServerManager) Start() {
	m.createDocker()
	m.createNetwork()
}
