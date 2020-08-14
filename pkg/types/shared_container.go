package types

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kardianos/service"
	"os"
	"strings"
)

// SharedContainer shares resources between the Daemon and its subordinates
type SharedContainer struct {
	Options *Options
	Logger  service.Logger

	Redis *redis.Client
}

// RedisSubscribe Returns a PubSub subscribed to the specified channels
func (c *SharedContainer) RedisSubscribe(channels ...string) *redis.PubSub {
	cp := redis.NewClient(c.Redis.Options())
	return cp.Subscribe(context.Background(), channels...)
}

// RedK returns a correctly prefixed Redis key to save data
func (c *SharedContainer) RedK(parts ...string) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("guardian.%s.%s", hostname, strings.Join(parts, "."))
}
