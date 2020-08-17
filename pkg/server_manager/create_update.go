package server_manager

import (
	"context"
	"encoding/json"
	"fmt"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/vanillaverse/guardian/pkg/types"
	"io/ioutil"
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
			return fmt.Errorf("failed while retrieving image: %v", err)
		}
	}

	if server.ImageID == "" {
		return fmt.Errorf("failed to create server: image not found")
	}

	networkMode := container.NetworkMode("host")
	mounts := []mount.Mount{}
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
			Mounts:       mounts,
			NetworkMode:  networkMode,
		},
		&network.NetworkingConfig{
			EndpointsConfig: endpointsConfig,
		},
		server.Name,
	)
	if err != nil {
		return fmt.Errorf("error while creating container: %v", err)
	}

	err = m.docker.ContainerStart(
		context.Background(),
		createRes.ID,
		dockerTypes.ContainerStartOptions{},
	)
	if err != nil {
		return fmt.Errorf("error while starting container: %v", err)
	}

	id := createRes.ID[:12]
	srvJson, _ := json.Marshal(server)

	m.Redis.LPush(context.Background(), m.RedisKey("managed"), id)
	m.Redis.Set(context.Background(), m.RedisKey("server", server.Name), srvJson, 0)

	m.Info("started server %s with container ID %s", server.Name, id)
	return nil
}

func (m *ServerManager) Update(server *types.Server) error {
	if server.Name == "" {
		return fmt.Errorf("failed to update server: name empty")
	}

	var err error
	server.ImageID, err = m.pullImage(server)
	if err != nil {
		return fmt.Errorf("failed while retrieving image: %v", err)
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
		return fmt.Errorf("error while stopping container: %v", err)
	}

	// Attempt to remove container (in non-persistent it should not exist, this is a safeguard)
	_ = m.docker.ContainerRemove(
		context.Background(),
		server.Name,
		dockerTypes.ContainerRemoveOptions{},
	)

	// Wrap recreation error
	if err := m.Create(server); err != nil {
		return fmt.Errorf("error while recreating server: %v", err)
	}

	return nil
}

func (m *ServerManager) pullImage(server *types.Server) (string, error) {
	imgName := fmt.Sprintf(m.Options.Docker.ImageFormat, server.Flavor, server.Version)
	pull, err := m.docker.ImagePull(
		context.Background(),
		imgName,
		dockerTypes.ImagePullOptions{RegistryAuth: m.credentials()},
	)
	if err != nil {
		return "", err
	}

	m.Info("pulling image for server %s: %s", server.Name, imgName)
	_, _ = ioutil.ReadAll(pull)
	_ = pull.Close()

	img, _, err := m.docker.ImageInspectWithRaw(context.Background(), imgName)
	if err != nil {
		return "", err
	}

	return img.ID, nil
}
