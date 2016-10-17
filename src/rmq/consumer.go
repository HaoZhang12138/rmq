package rmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"handler"
)

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	done       chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queue, key string) (*Consumer, error) {
	c := &Consumer{
		connection: nil,
		channel:    nil,
		done:       make(chan error),
	}

	var err error

	log.Printf("Connecting to %s", amqpURI)
	c.connection, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("hannel: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	q, err := c.channel.QueueDeclare(
		queue, // name of the queue
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		q.Name,    // name of the queue
		key,      // routingKey
		exchange, // sourceExchange
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	deliveries, err := c.channel.Consume(
		q.Name, // name
		 "", // consumerTag,
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	//if err := c.channel.Cancel(c.connection., true); err != nil {
	//	return fmt.Errorf("Consumer cancel failed: %s", err)
	//}
	log.Printf("shutdown")
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		err := handler.Handleafterconsume(d.Body)
		if err != nil {
			done <- err
			return
		}
	}
	log.Printf("Handle: deliveries channel closed")
	done <- nil
}