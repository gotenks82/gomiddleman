package middleman

import (
	. "GoMiddleMan/src/models"
	"github.com/google/uuid"
	"github.com/teivah/gosiris/gosiris"
	"log"
)

type userActor struct {
	gosiris.Actor
	user User
	actorRef gosiris.ActorRefInterface
}

func getMiddleManActorRef() gosiris.ActorRefInterface {
	actorRef, _ := gosiris.ActorSystem().ActorOf(middlemanActorName)
	return actorRef
}

func getUserActorRef(name string) gosiris.ActorRefInterface {
	actorRef, err := gosiris.ActorSystem().ActorOf(name)

	if err != nil {
		actorRef = addUserActor(name)
	}

	return actorRef
}

func addUserActor(name string) gosiris.ActorRefInterface {
	user := User{
		Id:        name,
		Interests: make([]Interest, 0),
	}
	actor := userActor{
		Actor: gosiris.Actor{},
		user:  user,
	}
	actor.addHandlers()

	_ = gosiris.ActorSystem().RegisterActor(name, &actor, nil)
	actorRef, err := gosiris.ActorSystem().ActorOf(name)

	if err != nil {
		log.Fatal("Could not start user actor for name", name)
	} else {
		actor.actorRef = actorRef
		_ = actorRef.Tell(gosiris.EmptyContext, "init", nil, actorRef)
	}
	return actorRef
}

func (actor *userActor) addHandlers() {
	actor.React("init", func(ctx gosiris.Context) {
		log.Print("user actor initialized: ", actor.Name())
	})

	actor.React(addInterestMsg, func(ctx gosiris.Context) {
		userInterest := ctx.Data.(UserInterest)
		log.Print(actor.Name(), ", Received interest: ", userInterest)
		actor.addInterest(userInterest)
	})

	actor.React(tradeMsg, func(ctx gosiris.Context) {
		trade := ctx.Data.(TradeOpportunity)
		log.Print(actor.Name(), ", Received trade: ", trade)
		if trade.IsComplete() {
			log.Printf("Trade %s is complete!", trade.Id)
		} else if !trade.IsUserInvolved(actor.user.Id) {
			actor.forwardTrade(trade)
		}
	})
}

func (actor *userActor) addInterest(userInterest UserInterest) {
	actor.user.AddInterest(userInterest.Interest)
	actor.sendToUser(msgToUser{
		dest:    userInterest.Interest.ItemUserid,
		msgType: tradeMsg,
		data: TradeOpportunity{
			Id:    uuid.New().String(),
			RootUserId: userInterest.UserId,
			Steps: []UserInterest{userInterest},
		},
	})
}

func (actor *userActor) forwardTrade(opportunity TradeOpportunity) {
	for _, interest := range actor.user.Interests {
		if !opportunity.WasReceivedBy(interest.ItemUserid) {
			actor.sendToUser(msgToUser{
				dest:    interest.ItemUserid,
				msgType: tradeMsg,
				data: opportunity.CopyWithStep(UserInterest{
					UserId:   actor.user.Id,
					Interest: interest,
				}),
			})
		}
	}
}
