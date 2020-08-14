package server_manager

import (
	"github.com/goccy/go-yaml"
	"github.com/vanillaverse/guardian/pkg/types"
	"io/ioutil"
)

func (m *ServerManager) PopulateServers() {
	serversYml, err := ioutil.ReadFile(m.Options.ServersFile)
	if err != nil {
		m.Warn("Failed to read server.yml, skipping: %v", err)
		return
	}

	var servers []*types.Server
	err = yaml.Unmarshal(serversYml, &servers)
	if err != nil {
		m.Warn("Failed to parse server.yml, skipping: %v", err)
		return
	}

	for _, srv := range servers {
		if err := m.Create(srv); err != nil {
			m.Error("failed while provisioning %s: %v", srv.Name, err)
		}
	}
}
