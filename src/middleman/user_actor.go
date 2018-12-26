package middleman

import (
	. "GoMiddleMan/src/models"
	"github.com/google/uuid"
	"github.com/gotenks82/gosiris/gosiris"
	"log"
)

type userActor struct {
	gosiris.Actor
	user     User
	actorRef gosiris.ActorRefInterface
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
		Id:            name,
		Interests:     make([]Interest, 0),
		Trades:        make([]TradeOpportunity, 0),
		Notifications: make([]string, 0),
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
		_ = actorRef.Tell(gosiris.EmptyContext, initMsg, nil, actorRef)
	}
	return actorRef
}

func (actor *userActor) addHandlers() {
	actor.React(initMsg, func(ctx gosiris.Context) {
		log.Print("user actor initialized: ", actor.Name())
	}).React(addInterestMsg, func(ctx gosiris.Context) {
		userInterest := ctx.Data.(UserInterest)
		log.Print(actor.Name(), ", Received interest: ", userInterest)
		actor.addInterest(userInterest)
	}).React(tradeOpportunityMsg, func(ctx gosiris.Context) {
		trade := ctx.Data.(TradeOpportunity)
		log.Print(actor.Name(), ", Received trade: ", trade)
		if trade.IsComplete() {
			log.Printf("Trade %s is complete!", trade.Id)
			actor.createTrade(trade)
		} else if !trade.IsUserInvolved(actor.user.Id) {
			actor.forwardTrade(trade)
		}
	}).React(storeTradeMsg, func(ctx gosiris.Context) {
		trade := ctx.Data.(TradeOpportunity)
		actor.storeTrade(trade)
	}).React(getNotificationsMsg, func(ctx gosiris.Context) {
		replyMsgType := ctx.Data.(question).replyMsgType
		notifications := actor.user.GetAndResetNotifications()
		_ = ctx.Sender.Tell(gosiris.EmptyContext, replyMsgType, notifications, ctx.Self)
	})
}

func (actor *userActor) addInterest(userInterest UserInterest) {
	actor.user.AddInterest(userInterest.Interest)
	actor.sendToUser(msgToUser{
		dest:    userInterest.Interest.ItemUserid,
		msgType: tradeOpportunityMsg,
		data: TradeOpportunity{
			Id:         uuid.New().String(),
			RootUserId: userInterest.UserId,
			Steps:      []UserInterest{userInterest},
		},
	})
}

func (actor *userActor) createTrade(trade TradeOpportunity) {
	actor.sendToTrade(msgToTrade{
		dest:    trade.Id,
		msgType: createTradeMsg,
		data:    trade,
	})
}

func (actor *userActor) forwardTrade(trade TradeOpportunity) {
	for _, interest := range actor.user.Interests {
		if !trade.WasReceivedBy(interest.ItemUserid) {
			actor.sendToUser(msgToUser{
				dest:    interest.ItemUserid,
				msgType: tradeOpportunityMsg,
				data: trade.CopyWithStep(UserInterest{
					UserId:   actor.user.Id,
					Interest: interest,
				}),
			})
		}
	}
}

func (actor *userActor) storeTrade(trade TradeOpportunity) {
	actor.user.AddTrade(trade)
	log.Printf("User %s stored Trade %s", actor.user.Id, trade.Id)
	actor.user.AddNotification("You have a new trade opportunity!")
}

func (actor *userActor) sendToUser(msg msgToUser) {
	sendMessage(sendToUserMsg, msg, actor.actorRef)
}

func (actor *userActor) sendToTrade(msg msgToTrade) {
	sendMessage(sendToTradeMsg, msg, actor.actorRef)
}
