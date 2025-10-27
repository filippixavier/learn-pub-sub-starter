package pubsub

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Transient SimpleQueueType = iota
	Durable
)

func PublishJSON[T any](ch *amqp091.Channel, exchange, key string, val T) error {
	body, err := json.Marshal(val)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	return err
}

func DeclareAndBind(
	conn *amqp091.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp091.Channel, amqp091.Queue, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	queue, err := ch.QueueDeclare(queueName, queueType == Durable, queueType == Transient, queueType == Transient, false, nil)

	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)

	if err != nil {
		return nil, amqp091.Queue{}, err
	}

	return ch, queue, nil
}
