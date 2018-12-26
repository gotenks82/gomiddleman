package middleman

import (
	"github.com/gotenks82/gosiris/gosiris"
	"log"
)

type (
	middleManActor struct {
		gosiris.Actor
		actorRef gosiris.ActorRefInterface
	}

	msgToUser struct {
		dest    string
		msgType string
		data    interface{}
	}

	askToUser struct {
		dest         string
		question     question
	}

	question struct {
		msgType		 string
		data         interface{}
		replyMsgType string
	}

	msgToTrade struct {
		dest    string
		msgType string
		data    interface{}
	}
)

func addMiddlemanActor() (gosiris.ActorRefInterface, error) {
	manager := middleManActor{}
	manager.addHandlers()

	_ = gosiris.ActorSystem().RegisterActor(middlemanActorName, &manager, nil)
	actorRef, err := gosiris.ActorSystem().ActorOf(middlemanActorName)

	manager.actorRef = actorRef

	if err != nil {
		log.Fatal("Could not start middleman actor")
	} else {
		_ = actorRef.Tell(gosiris.EmptyContext, "init", nil, actorRef)
	}
	return actorRef, err
}

func (actor *middleManActor) addHandlers() {
	actor.React("init", func(ctx gosiris.Context) {
		log.Print("middleMan actor initialized")
	}).React(sendToUserMsg, func(ctx gosiris.Context) {
		sendToUser := ctx.Data.(msgToUser)
		log.Print(middlemanActorName, ", Received msgToUser: ", sendToUser)

		userRef := getUserActorRef(sendToUser.dest)
		err := userRef.Tell(ctx, sendToUser.msgType, sendToUser.data, ctx.Sender)

		if err != nil {
			log.Print("Could not forward message to user: ", sendToUser.dest)
		}
	}).React(sendToTradeMsg, func(ctx gosiris.Context) {
		sendToTrade := ctx.Data.(msgToTrade)
		log.Print(middlemanActorName, ", Received msgToTrade: ", sendToTrade)

		tradeRef := getTradeActorRef(sendToTrade.dest)
		err := tradeRef.Tell(ctx, sendToTrade.msgType, sendToTrade.data, ctx.Sender)

		if err != nil {
			log.Print("Could not forward message to trade: ", sendToTrade.dest)
		}
	}).React(askToUserMsg, func(ctx gosiris.Context) {
		askToUser := ctx.Data.(askToUser)
		log.Print(middlemanActorName, ", Received msgToUser: ", askToUser)

		userRef := getUserActorRef(askToUser.dest)
		err := userRef.Tell(ctx, askToUser.question.msgType, askToUser.question, ctx.Sender)

		if err != nil {
			log.Print("Could not forward message to user: ", askToUser.dest)
		}
	})
}
