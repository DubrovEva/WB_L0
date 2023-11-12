package nats_sub

import (
	"encoding/json"
	"fmt"
	"github.com/DubrovEva/WB_L0/internal/db"
	"github.com/DubrovEva/WB_L0/internal/entities"
	"github.com/nats-io/stan.go"
	"log"
)

type MessageProcessor struct {
	conn stan.Conn
	db   *db.Database
	sub  stan.Subscription
}

func Start(cluster, subClient, channel, URL string, db *db.Database) (*MessageProcessor, error) {
	var mp MessageProcessor
	mp.db = db
	sc, err := stan.Connect(cluster, subClient, stan.NatsURL(URL))
	if err != nil {
		return nil, fmt.Errorf("can't connect: %v", err)
	}
	mp.conn = sc
	log.Println("Stan connection established")

	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		err := mp.process(msg)
		if err != nil {
			log.Printf("Can't process message: %v", err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("can't subscribe: %v", err)
	}
	mp.sub = sub
	log.Println("Subscribed successfully")

	return &mp, nil
}

func (mp *MessageProcessor) Stop() error {
	err := mp.conn.Close()
	if err != nil {
		return err
	}

	err = mp.sub.Unsubscribe()
	if err != nil {
		return err
	}
	return nil
}

func (mp *MessageProcessor) process(m *stan.Msg) error {
	log.Println("Got new message")

	order := entities.Order{}
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		return fmt.Errorf("can't unmarshal message: %v", err)
	}

	err = mp.db.SaveOrder(&order)
	if err != nil {
		return fmt.Errorf("can't save order: %v", err)
	}
	log.Println("Order was saved")

	return nil
}
