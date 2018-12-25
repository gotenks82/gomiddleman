package middleman

import (
	"github.com/teivah/gosiris/gosiris"
	"log"
)

type (
	middleManActor struct {
		gosiris.Actor
	}

	msgToUser struct {
		dest string
		msgType string
		data interface{}
	}
)



func addMiddlemanActor() (gosiris.ActorRefInterface, error) {
	manager := middleManActor{}
	manager.addHandlers()

	_ = gosiris.ActorSystem().RegisterActor(middlemanActorName, &manager, nil)
	actorRef, err := gosiris.ActorSystem().ActorOf(middlemanActorName)

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
	})

	actor.React(sendToUserMsg, func(ctx gosiris.Context) {
		sendToUser := ctx.Data.(msgToUser)
		log.Print(middlemanActorName, ", Received msgToUser: ", sendToUser)

		userRef := getUserActorRef(sendToUser.dest)
		err := userRef.Tell(ctx, sendToUser.msgType, sendToUser.data, ctx.Sender)

		if err != nil {
			log.Print("Could not forward message to user: ", sendToUser.dest)
		}
	})
}


