package middleman

import (
	. "GoMiddleMan/src/models"
	"errors"
	"github.com/gotenks82/gosiris/gosiris"
	"log"
	"sync"
	"time"
)

const middlemanActorName = "middleman"
const defaultTimeout = 3 * time.Second

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

func getMiddleManActorRef() gosiris.ActorRefInterface {
	actorRef, err := gosiris.ActorSystem().ActorOf(middlemanActorName)
	if err != nil {
		log.Fatal("Could not find middleman actor")
	}
	return actorRef
}

func (m *middleMan) AddInterest(userId string, i Interest) {
	actorRef := getMiddleManActorRef()
	if actorRef != nil {
		msgToUser := msgToUser{
			dest:    userId,
			msgType: addInterestMsg,
			data: UserInterest{
				UserId:   userId,
				Interest: i,
			},
		}
		_ = actorRef.Tell(gosiris.EmptyContext, sendToUserMsg, msgToUser, actorRef)
	}
}

func (m *middleMan) GetNotifications(userId string) []string {
	var notifications []string
	reply, _ := askQuestionToUser(userId, getNotificationsMsg, nil)
	if reply != nil {
		notifications = reply.([]string)
	}
	return notifications
}

func askQuestionToUser(userId string, msgType string, data interface{}) (interface{}, error) {
	actorRef := getMiddleManActorRef()

	if actorRef != nil {
		askToUser := askToUser{
			dest: userId,
			question: question{
				msgType:      msgType,
				replyMsgType: askToUserMsg,
				data:         data,
			},
		}
		return actorRef.Ask(gosiris.EmptyContext, askToUserMsg, askToUser, defaultTimeout)
	}

	return nil, errors.New("could not get middleman actor")
}

func sendMessage(msgType string, msg interface{}, src gosiris.ActorRefInterface) {
	actorRef := getMiddleManActorRef()
	if actorRef != nil {
		err := actorRef.Tell(gosiris.EmptyContext, msgType, msg, src)

		if err != nil {
			log.Print(err.Error())
		}
	}
}
