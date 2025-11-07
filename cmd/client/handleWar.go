package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(war gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(war)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits,
			gamelogic.WarOutcomeOpponentWon,
			gamelogic.WarOutcomeYouWon,
			gamelogic.WarOutcomeDraw:
			return pubsub.Ack
		default:
			fmt.Println("an error as occured when trying to process war")
			return pubsub.NackDiscard
		}
	}
}
