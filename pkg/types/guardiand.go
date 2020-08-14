package types

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kardianos/service"
	"log"
	"os"
	"strings"
)

// GuardianD shares resources between the Daemon and its subordinates
type GuardianD struct {
	Options *Options
	Logger  service.Logger

	Redis *redis.Client
}

// RedisSubscribe Returns a PubSub subscribed to the specified channels
func (g *GuardianD) RedisSubscribe(channels ...string) *redis.PubSub {
	cp := redis.NewClient(g.Redis.Options())
	return cp.Subscribe(context.Background(), channels...)
}

// RedisKey returns a correctly prefixed Redis key to save data
func (g *GuardianD) RedisKey(parts ...string) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("guardian.%s.%s", hostname, strings.Join(parts, "."))
}

// Info creates an info level log
func (g *GuardianD) Info(format string, params ...interface{}) {
	if err := g.Logger.Infof(format, params...); err != nil {
		log.Fatalf("error while logging: %v", err)
	}
}

// Error creates an error level log
func (g *GuardianD) Error(format string, params ...interface{}) {
	if err := g.Logger.Errorf(format, params...); err != nil {
		log.Fatalf("error while logging: %v", err)
	}
}

// Warn creates an warn level log
func (g *GuardianD) Warn(format string, params ...interface{}) {
	if err := g.Logger.Warningf(format, params...); err != nil {
		log.Fatalf("error while logging: %v", err)
	}
}
