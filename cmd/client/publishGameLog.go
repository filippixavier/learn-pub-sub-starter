package main

import (
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(ch *amqp.Channel, username, initiator, message string) error {
	gameLog := routing.GameLog{
		CurrentTime: time.Now(),
		Username:    username,
		Message:     message,
	}
	return pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+initiator, gameLog)
}
