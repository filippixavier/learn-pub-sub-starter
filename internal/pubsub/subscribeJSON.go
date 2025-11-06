package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)

	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		for msg := range deliveryCh {
			var body T
			err := json.Unmarshal(msg.Body, &body)

			if err != nil {
				fmt.Println(err)
			} else {
				status := handler(body)

				switch status {
				case Ack:
					fmt.Println("Msg ACK")
					msg.Ack(false)
				case NackRequeue:
					fmt.Println("Msg Nack, requeue")
					msg.Nack(false, true)
				case NackDiscard:
					fmt.Println("Msg Nack, no requeue")
					msg.Nack(false, false)
				}
			}
		}
	}()

	return nil
}
