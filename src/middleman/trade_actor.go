package middleman

import (
	"GoMiddleMan/src/models"
	"github.com/gotenks82/gosiris/gosiris"
	"log"
)

type tradeActor struct {
	gosiris.Actor
	trade    models.TradeOpportunity
	actorRef gosiris.ActorRefInterface
}

func getTradeActorRef(name string) gosiris.ActorRefInterface {
	actorRef, err := gosiris.ActorSystem().ActorOf(name)

	if err != nil {
		actorRef = addTradeActor(name)
	}

	return actorRef
}

func addTradeActor(name string) gosiris.ActorRefInterface {
	actor := tradeActor{
		Actor: gosiris.Actor{},
	}
	actor.addHandlers()

	_ = gosiris.ActorSystem().RegisterActor(name, &actor, nil)
	actorRef, err := gosiris.ActorSystem().ActorOf(name)

	if err != nil {
		log.Fatal("Could not start trade actor for name", name)
	} else {
		actor.actorRef = actorRef
		_ = actorRef.Tell(gosiris.EmptyContext, initMsg, nil, actorRef)
	}
	return actorRef
}

func (actor *tradeActor) addHandlers() {
	actor.React(initMsg, func(ctx gosiris.Context) {
		log.Print("trade actor initialized: ", actor.Name())
	}).React(createTradeMsg, func(ctx gosiris.Context) {
		trade := ctx.Data.(models.TradeOpportunity)
		actor.trade = trade
		actor.notifyCreation()
	})
}

func (actor *tradeActor) notifyCreation() {
	for _, step := range actor.trade.Steps {
		actor.sendToUser(msgToUser{
			dest: step.UserId,
			msgType: storeTradeMsg,
			data: actor.trade,
		})
	}
}

func (actor *tradeActor) sendToUser(msg msgToUser) {
	sendMessage(sendToUserMsg, msg, actor.actorRef)
}