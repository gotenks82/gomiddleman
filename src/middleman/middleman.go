package middleman

import (
	. "GoMiddleMan/src/models"
	"github.com/teivah/gosiris/gosiris"
	"log"
	"sync"
)

const middlemanActorName = "middleman"

type (
	middleMan struct {
		middleManActorRef gosiris.ActorRefInterface
	}
)

var instance *middleMan
var once sync.Once

func GetInstance() *middleMan {
	once.Do(func() {
		err := gosiris.InitActorSystem(gosiris.SystemOptions{
			ActorSystemName: "ActorSystem",
		})
		if err == nil {
			instance = &middleMan{}
			instance.middleManActorRef, err = addMiddlemanActor()
		}
		if err != nil {
			log.Fatal(err.Error())
		}
	})
	return instance
}

func (m *middleMan) Shutdown() {
	log.Print("Shutting down actor system...")
	_ = gosiris.CloseActorSystem()
}

func (m *middleMan) AddInterest(userId string, i Interest) {
	actorRef, err := gosiris.ActorSystem().ActorOf(middlemanActorName)
	if err != nil {
		log.Fatal("Could not find middleman actor")
	} else {
		msgToUser := msgToUser{
			dest: userId,
			msgType: addInterestMsg,
			data: UserInterest{
				UserId:   userId,
				Interest: i,
			},
		}

		_ = actorRef.Tell(gosiris.EmptyContext, sendToUserMsg, msgToUser, actorRef)

	}
}

func (actor *userActor) sendToUser(msg msgToUser) {
	err := getMiddleManActorRef().Tell(gosiris.EmptyContext, sendToUserMsg, msg, actor.actorRef)

	if err != nil {
		log.Print(err.Error())
	}
}