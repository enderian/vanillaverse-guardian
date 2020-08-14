package server_manager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
)

const (
	DockerNetwork         = "vanillaverse"
	DockerServerHashLabel = "guardian/server/hash"
)

func (m *ServerManager) portMap(server *types.Server) nat.PortMap {
	_, mp, err := nat.ParsePortSpecs(server.Ports)
	if err != nil {
		log.Printf("error while parsing ports for server %v: %v", server.Name, err)
	}
	return mp
}

func (m *ServerManager) credentials() string {
	js, _ := json.Marshal(m.Options.DockerCredentials)
	return base64.StdEncoding.EncodeToString(js)
}

func (m *ServerManager) createNetwork() {
	_, _ = m.docker.NetworkCreate(context.Background(), DockerNetwork, dockerTypes.NetworkCreate{
		CheckDuplicate: true,
	})
}
