package rmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Producer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	done       chan error
}

func NewProducer(amqpURI, exchange, exchangeType string) (*Producer, error) {
	p := &Producer{
		connection: nil,
		channel:    nil,
		done:       make(chan error),
	}

	var err error

	log.Printf("Connecting to %s", amqpURI)
	p.connection, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: ", err)
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: ", err)
	}

	if err := p.channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	return p, nil
}

func (p *Producer) Publish(exchange string, routingKey string, body []byte) error {

	if err := p.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			Body:            body,
			DeliveryMode:    amqp.Persistent, // 1=non-persistent, 2=persistent
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: ", err)
	}

	return nil
}