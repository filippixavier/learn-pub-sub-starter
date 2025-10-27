package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
)

const CONNECTIONSTRING string = "amqp://guest:guest@localhost:5672/"

func main() {
	connexion, err := amqp091.Dial(CONNECTIONSTRING)

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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Server gracefully stopped with %v\n", sig)
}
