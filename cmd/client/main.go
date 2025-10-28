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

	pubsub.DeclareAndBind(connexion, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, username), routing.PauseKey, pubsub.Transient)

	gamestate := gamelogic.NewGameState(username)

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
			_, err := gamestate.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Move worked!")
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
