package server_manager

import (
	"context"
	"encoding/json"
	"fmt"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/vanillaverse/guardian/pkg/types"
	"io/ioutil"
	"log"
	"time"
)

func (m *ServerManager) Create(server *types.Server) error {
	if server.Name == "" {
		return fmt.Errorf("failed to create server: name empty")
	}

	if server.ImageID == "" {
		var err error
		server.ImageID, err = m.pullImage(server)
		if err != nil {
			return fmt.Errorf("failed retrieve image for server %s: %v", server.Name, err)
		}
	}

	if server.ImageID == "" {
		return fmt.Errorf("failed to create server: image not found")
	}

	networkMode := container.NetworkMode("host")
	endpointsConfig := map[string]*network.EndpointSettings{}

	if !server.OnHost {
		networkMode = "default"
		endpointsConfig[DockerNetwork] = &network.EndpointSettings{
			NetworkID: DockerNetwork,
			Aliases:   []string{server.Name},
		}
	}

	createRes, err := m.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: server.ImageID,
			Labels: map[string]string{
				DockerServerHashLabel: server.Hash(),
			},
		},
		&container.HostConfig{
			AutoRemove:   !server.Persist,
			PortBindings: m.portMap(server),
			// TODO: Mounts
			NetworkMode: networkMode,
		},
		&network.NetworkingConfig{
			EndpointsConfig: endpointsConfig,
		},
		server.Name,
	)
	if err != nil {
		return fmt.Errorf("error while creating server %s: %v", server.Name, err)
	}

	err = m.docker.ContainerStart(
		context.Background(),
		createRes.ID,
		dockerTypes.ContainerStartOptions{},
	)
	if err != nil {
		return fmt.Errorf("error while starting server %s: %v", server.Name, err)
	}

	id := createRes.ID[:12]
	srvJson, _ := json.Marshal(server)

	m.Redis.LPush(context.Background(), m.RedisKey("managed"), id)
	m.Redis.Set(context.Background(), m.RedisKey("server", server.Name), srvJson, 0)

	log.Printf("started server %s with container ID %s", server.Name, id)
	return nil
}

func (m *ServerManager) Update(server *types.Server) error {
	if server.Name == "" {
		return fmt.Errorf("failed to update server: name empty")
	}

	var err error
	server.ImageID, err = m.pullImage(server)
	if err != nil {
		return fmt.Errorf("failed retrieve image for server %s: %v", server.Name, err)
	}

	// Determine if server requires update
	cont, err := m.docker.ContainerInspect(context.Background(), server.Name)
	if err == nil && cont.Config.Labels[DockerServerHashLabel] == server.Hash() {
		return fmt.Errorf("%s does not require an update, stopping process", server.Name)
	}

	// Stop (this blocks until stopped) the server
	duration := 65 * time.Second
	err = m.docker.ContainerStop(context.Background(), server.Name, &duration)
	if err != nil {
		return fmt.Errorf("error while stopping server %s: %v", server.Name, err)
	}

	// Attempt to remove container (in non-persistent it should not exist, this is a safeguard)
	_ = m.docker.ContainerRemove(
		context.Background(),
		server.Name,
		dockerTypes.ContainerRemoveOptions{},
	)

	// Wrap recreation error
	if err := m.Create(server); err != nil {
		return fmt.Errorf("error while recreating server %s: %v", server.Name, err)
	}

	return nil
}

func (m *ServerManager) pullImage(server *types.Server) (string, error) {
	pull, err := m.docker.ImagePull(
		context.Background(),
		fmt.Sprintf(m.Options.DockerImageFormat, server.Flavor, server.Version),
		dockerTypes.ImagePullOptions{RegistryAuth: m.credentials()},
	)
	if err != nil {
		return "", fmt.Errorf("error while pulling image for server %s: %v", server.Name, err)
	}

	m.Info("pulling image for server %s", server.Name)
	_, _ = ioutil.ReadAll(pull)
	_ = pull.Close()

	img, _, err := m.docker.ImageInspectWithRaw(context.Background(), server.Name)
	if err != nil {
		return "", fmt.Errorf("error while inspecting image for server %s: %v", server.Name, err)
	}

	return img.ID, nil
}
