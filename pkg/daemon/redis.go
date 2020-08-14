package daemon

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/vanillaverse/guardian/pkg/types"
	"log"
)

const (
	RedisServerCreate = "server.create"
)

func (d *Daemon) initializeRedis() {
	opts, err := redis.ParseURL(d.Options.RedisURL)
	if err != nil {
		log.Fatalf("failed to parse redis URL (%s): %v", d.Options.RedisURL, err)
	}
	d.Redis = redis.NewClient(opts)
	log.Printf("initialized redis client to %s", opts.Addr)
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
			d.srvManager.Create(srv)
		}
	}
}
