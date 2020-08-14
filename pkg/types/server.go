package types

import (
	"crypto/sha256"
	"fmt"
)

// Server represents a Minecraft server instance
type Server struct {
	Name    string `json:"name"`
	Flavor  string `json:"flavor"`
	Version string `json:"version"`

	Ports    []string `json:"ports"`
	Volumes  []string `json:"volumes"`
	Mappings []string `json:"mappings"`

	OnHost  bool `json:"on_host"`
	Persist bool `json:"persist"`

	// Never set this variable by hand
	// It should only be handled by the daemon
	ImageID string `json:"image_id"`
}

// Hash is the digest of the server struct
func (s Server) Hash() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", s)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
