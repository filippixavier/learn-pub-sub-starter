package main

import (
	"fmt"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CONNECTIONSTRING string = "amqp://guest:guest@localhost:5672/"

func main() {
	connexion, err := amqp.Dial(CONNECTIONSTRING)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer connexion.Close()
	fmt.Printf("Connection to %s successful!\n", CONNECTIONSTRING)

	ch, err := connexion.Channel()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pubsub.SubscribeGob(connexion, routing.ExchangePerilTopic, routing.GameLogSlug, fmt.Sprintf("%s.*", routing.GameLogSlug), pubsub.Durable, handlerLogs())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("Sending pause message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			fmt.Println("Sending resume message")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			fmt.Println("Exiting server...")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
