package queuedeclare

import (
	"encoding/json"
	"fmt"

	"github.com/setval/infr/rabbitmq/content"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	d amqp.Delivery
}

func RunDelivery(host, name string, chanEvent chan Delivery) error {
	conn, err := amqp.Dial(host)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open a channel %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("declare a queue: %w", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("set QoS: %w", err)
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
		return fmt.Errorf("register a consumer: %w", err)
	}

	for d := range msgs {
		chanEvent <- Delivery{d}
	}

	return nil
}

func (d *Delivery) Ack() {
	d.d.Ack(false)
}

func (d *Delivery) Content() (content.Content, error) {
	var c content.Content
	return c, json.Unmarshal(d.d.Body, &c)
}
