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

func (m *ServerManager) Create(server *types.Server) {
	if server.Name == "" {
		log.Println("failed to create server: name empty")
		return
	}

	if server.ImageID == "" {
		server.ImageID = m.pullImage(server)
	}
	if server.ImageID == "" {
		log.Println("failed to create server: image not found")
		return
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
		log.Printf("error while creating server %s: %v", server.Name, err)
		return
	}

	err = m.docker.ContainerStart(
		context.Background(),
		createRes.ID,
		dockerTypes.ContainerStartOptions{},
	)
	if err != nil {
		log.Printf("error while starting server %s: %v", server.Name, err)
		return
	}

	id := createRes.ID[:12]
	srvJson, _ := json.Marshal(server)

	m.Redis.LPush(context.Background(), m.RedK("managed"), id)
	m.Redis.Set(context.Background(), m.RedK("server", server.Name), srvJson, 0)

	log.Printf("started server %s with container ID %s", server.Name, id)
}

func (m *ServerManager) Update(server *types.Server) {
	if server.Name == "" {
		log.Println("failed to update server: name empty")
		return
	}

	server.ImageID = m.pullImage(server)
	if server.ImageID == "" {
		log.Println("failed to create server: image not found")
		return
	}

	// Determine if server requires update
	cont, err := m.docker.ContainerInspect(context.Background(), server.Name)
	if err == nil && cont.Config.Labels[DockerServerHashLabel] == server.Hash() {
		log.Printf("%s does not require an update, stopping process", server.Name)
		return
	}

	// Stop (this blocks until stopped) the server
	duration := 65 * time.Second
	err = m.docker.ContainerStop(context.Background(), server.Name, &duration)
	if err != nil {
		log.Printf("error while stopping server %s: %v", server.Name, err)
		return
	}

	// Attempt to remove container (in non-persistent it should not exist, this is a safeguard)
	_ = m.docker.ContainerRemove(
		context.Background(),
		server.Name,
		dockerTypes.ContainerRemoveOptions{},
	)

	m.Create(server)
}

func (m *ServerManager) pullImage(server *types.Server) string {
	pull, err := m.docker.ImagePull(
		context.Background(),
		fmt.Sprintf(m.Options.DockerImageFormat, server.Flavor, server.Version),
		dockerTypes.ImagePullOptions{RegistryAuth: m.credentials()},
	)
	if err != nil {
		log.Printf("error while pulling image for server %s: %v", server.Name, err)
		return ""
	}

	log.Printf("pulling image for server %s", server.Name)
	_, _ = ioutil.ReadAll(pull)
	_ = pull.Close()

	img, _, err := m.docker.ImageInspectWithRaw(context.Background(), server.Name)
	if err != nil {
		log.Printf("error while inspecting image for server %s: %v", server.Name, err)
		return ""
	}
	return img.ID
}
