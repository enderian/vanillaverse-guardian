package types

import (
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

type Options struct {
	ConfigFile  string
	ServersFile string

	RedisURL string `json:"redis_url"`

	DockerImageFormat string            `json:"docker_image_format"`
	DockerCredentials map[string]string `json:"docker_credentials"`
}

func (o *Options) ReadFromFiles() {
	configYml, err := ioutil.ReadFile(o.ConfigFile)
	if err != nil {
		log.Fatalf("error while reading config file: %v", err)
	}
	err = yaml.Unmarshal(configYml, o)
	if err != nil {
		log.Fatalf("error while parsing config file: %v", err)
	}
}
