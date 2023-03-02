package queuedeclare

import (
	"encoding/json"
	"fmt"

	"github.com/setval/infr/rabbitmq"
	"github.com/setval/infr/rabbitmq/content"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Delivery struct {
	d amqp.Delivery
}

func RunDelivery(host, name string, chanEvent chan Delivery, errCh chan error) error {
	conn, err := amqp.Dial(host)
	if err != nil {
		return fmt.Errorf("connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	notify := conn.NotifyClose(make(chan *amqp.Error))
	block := conn.NotifyBlocked(make(chan amqp.Blocking))

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

	for {
		select {
		case err = <-notify:
			errCh <- rabbitmq.ErrCloseConnecton
		case <-block:
			errCh <- rabbitmq.ErrBlockConnection
		case d := <-msgs:
			chanEvent <- Delivery{d}
		}
	}
}

func (d *Delivery) Ack() {
	d.d.Ack(false)
}

func (d *Delivery) Content() (content.Content, error) {
	var c content.Content
	return c, json.Unmarshal(d.d.Body, &c)
}
