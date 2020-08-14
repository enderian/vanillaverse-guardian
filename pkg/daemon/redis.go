package daemon

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
	"os"
)

const (
	RedisServerCreate = "server.create"
)

func (d *Daemon) initializeRedis() {
	opts, err := redis.ParseURL(d.Options.RedisURL)
	if err != nil {
		d.Error("failed to parse redis URL (%s): %v", d.Options.RedisURL, err)
		os.Exit(1)
	}
	d.Redis = redis.NewClient(opts)
	d.Info("initialized redis client to %s", opts.Addr)
	go d.listenRedis()
}

func (d *Daemon) listenRedis() {
	ps := d.RedisSubscribe(RedisServerCreate)
	_, err := ps.Receive(d.ctx)
	if err != nil {
		log.Fatalf("error while acquiring pubsub: %v", err)
	}

	ch := ps.Channel()
	for {
		msg := <-ch
		switch msg.Channel {
		case RedisServerCreate:
			srv := &types.Server{}
			_ = json.Unmarshal([]byte(msg.Payload), srv)
			_ = d.srvManager.Create(srv)
		}
	}
}
