package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// setupRabbit setups rabbit connection for saga purposes
func setupRabbit() (*amqp.Connection, error) {
	connection, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	if err != nil {
		return nil, fmt.Errorf("error on dialing to rabbit. err=%w", err)
	}

	return connection, nil
}
