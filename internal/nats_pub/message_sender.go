package nats_pub

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

var file1 = "internal/nats_pub/tests/first_test.json"
var file2 = "internal/nats_pub/tests/second_test.json"

type MessageSender struct {
	conn stan.Conn
}

func Start(cluster, client, channel, URL string) (*MessageSender, error) {
	var ms MessageSender
	conn, err := stan.Connect(cluster, client, stan.NatsURL(URL))
	if err != nil {
		return nil, err
	}
	ms.conn = conn
	log.Println("Stan connection established")

	source := file1
	testData, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	err = ms.conn.Publish(channel, testData)
	if err != nil {
		return nil, err
	}
	log.Printf("Published test data from file %s\n", source)

	return &ms, nil
}

func (ms *MessageSender) Stop() error {
	err := ms.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
