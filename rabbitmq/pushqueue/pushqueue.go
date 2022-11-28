package pushqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/setval/infr/rabbitmq/content"
)

type PushQueue struct {
	Q  *amqp.Queue
	Ch *amqp.Channel
}

func New(host, name string) (*PushQueue, error) {
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

	q, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("declare a queue: %w", err)
	}

	return &PushQueue{
		Q:  &q,
		Ch: ch,
	}, nil

}

func (p *PushQueue) Push(c content.Content) error {
	body, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal content: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.Ch.PublishWithContext(
		ctx,
		"",       // exchange
		p.Q.Name, // routing key
		false,    // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	return nil
}
