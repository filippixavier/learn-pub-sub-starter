package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
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
			buffer := bytes.NewBuffer(msg.Body)
			decoder := gob.NewDecoder(buffer)
			var body T

			err := decoder.Decode(&body)

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
