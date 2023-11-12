package main

import (
	"github.com/DubrovEva/WB_L0/internal/config"
	"github.com/DubrovEva/WB_L0/internal/nats_pub"
	"log"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	ms, err := nats_pub.Start(cfg.Nats.Cluster, cfg.Nats.PubClient, cfg.Nats.Channel, cfg.Nats.URL)
	if err != nil {
		log.Fatalf("Can't initialise publisher: %v", err)
	}
	defer func(ms *nats_pub.MessageSender) {
		err := ms.Stop()
		if err != nil {
			log.Printf("Could not close nats connection: %v", err)
		}
	}(ms)
}
