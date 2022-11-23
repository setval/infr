package recvdirect

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	d  *<-chan amqp.Delivery
	ch *amqp.Channel
}

func NewDelivery(host, exchange, key string) (*Delivery, error) {
	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, fmt.Errorf("connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open a channel %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("declare an exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("declare a queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		key,      // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("bind a queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("register a consumer: %w", err)
	}

	return &Delivery{d: &msgs, ch: ch}, nil
}

func (d *Delivery) Delivery() *<-chan amqp.Delivery {
	return d.d
}
