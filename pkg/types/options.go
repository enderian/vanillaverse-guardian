package types

import (
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
)

type dockerOptions struct {
	Host       string `json:"host"`
	VerifyHost bool   `json:"verify_host"`

	ImageFormat string            `json:"image_format"`
	Credentials map[string]string `json:"credentials"`
}

type Options struct {
	ConfigFile  string
	ServersFile string
	UnixSocket  string

	RedisURL string        `json:"redis_url"`
	Docker   dockerOptions `json:"docker"`
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
