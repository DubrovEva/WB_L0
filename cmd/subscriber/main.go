package main

import (
	"database/sql"
	"fmt"
	config "github.com/DubrovEva/WB_L0/internal/config"
	"github.com/DubrovEva/WB_L0/internal/db"
	httpserver "github.com/DubrovEva/WB_L0/internal/http-server"
	"github.com/DubrovEva/WB_L0/internal/nats_sub"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	fmt.Println(cfg.Nats.URL)

	database, err := db.InitializeDB(
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Host,
		cfg.DB.Port,
	)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("Could not close: %v", err)
		}
	}(database.Conn)

	mp, err := nats_sub.Start(cfg.Nats.Cluster, cfg.Nats.SubClient, cfg.Nats.Channel, cfg.Nats.URL, database)
	if err != nil {
		log.Fatalf("Could not initialize nats connection: %v", err)
	}
	defer func(mp *nats_sub.MessageProcessor) {
		err := mp.Stop()
		if err != nil {
			log.Printf("Could not close nats connection: %v", err)
		}
	}(mp)

	hs, err := httpserver.Start(database, cfg.HTTPServer.URL)
	if err != nil {
		log.Fatalf("Could not initialize http server: %v", err)
	}
	defer func(hs *httpserver.HttpServer) {
		err := hs.Stop()
		if err != nil {
			log.Printf("Could not stop http server: %v", err)
		}
	}(hs)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")

}
