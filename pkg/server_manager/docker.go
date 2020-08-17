package server_manager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/vanillaverse/guardian/pkg/types"
	"os"
)

const (
	DockerNetwork         = "vanillaverse"
	DockerServerHashLabel = "guardian/server/hash"
)

func (m *ServerManager) createDocker() {
	host := m.Options.Docker.Host
	if host == "" {
		host = client.DefaultDockerHost
	}

	cli, err := client.NewClient(host, client.DefaultVersion, nil, nil)
	if err != nil {
		m.Error("error while initializing docker client: %v", err)
		os.Exit(1)
	}

	m.Info("initialized docker client to %s", host)
	m.docker = cli
}

func (m *ServerManager) portMap(server *types.Server) nat.PortMap {
	_, mp, err := nat.ParsePortSpecs(server.Ports)
	if err != nil {
		m.Error("error while parsing ports for server %v: %v", server.Name, err)
	}
	return mp
}

func (m *ServerManager) credentials() string {
	js, _ := json.Marshal(m.Options.Docker.Credentials)
	return base64.StdEncoding.EncodeToString(js)
}

func (m *ServerManager) createNetwork() {
	_, _ = m.docker.NetworkCreate(context.Background(), DockerNetwork, dockerTypes.NetworkCreate{
		CheckDuplicate: true,
	})
}
