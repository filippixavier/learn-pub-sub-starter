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
		return
	}

	username, err := gamelogic.ClientWelcome()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	playerQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)

	gamestate := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(connexion, routing.ExchangePerilTopic, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gamestate))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	err = pubsub.SubscribeJSON(connexion, routing.ExchangePerilTopic, playerQueueName, routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerMove(gamestate, ch))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	err = pubsub.SubscribeJSON(connexion, routing.ExchangePerilTopic, "war", routing.WarRecognitionsPrefix+".*", pubsub.Durable, handlerWar(gamestate))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gamestate.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			move, err := gamestate.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			} else {
				err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, playerQueueName, move)

				if err == nil {
					fmt.Println("Move successfully broadcasted")
				} else {
					fmt.Println(err)
				}
			}
		case "status":
			gamestate.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
