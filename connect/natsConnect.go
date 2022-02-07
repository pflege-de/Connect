package connect

import (
	"log"

	"github.com/nats-io/nats.go"
)

func GetNatsConnection(host string, port string) (*nats.Conn, error) {
	nc, err := nats.Connect(host + ":" + port)
	if err != nil {
		log.Printf("Couldn't establish nats connection, please check if host is reachable and config values are right.\r\nError:%s", err)
		return nil, err
	}
	return nc, nil
}

func GetNatsEncoderFromConnection(n *nats.Conn) (*nats.EncodedConn, error) {
	ec, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
	if err != nil {
		log.Printf("Can't create Nats Encoder:\r\n %s", err)
		return nil, err
	}
	return ec, nil
}

func GetNatsEncoderFromSettings(host string, port string, n *nats.Conn) (*nats.EncodedConn, error) {
	ec, err := nats.NewEncodedConn(n, nats.JSON_ENCODER)
	if err != nil {
		log.Printf("Can't create Nats Encoder:\r\n %s", err)
		return nil, err
	}
	return ec, nil
}
