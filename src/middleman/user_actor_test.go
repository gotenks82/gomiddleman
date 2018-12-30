package middleman

import (
	. "GoMiddleMan/src/models"
	"github.com/gotenks82/gosiris/gosiris"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	instance := GetInstance()

	defer instance.Shutdown()
	os.Exit(m.Run())
}

func TestAddUserActor(t *testing.T) {
	t.Log("starting addUserActor test")

	userRef := addUserActor("testUser")

	if userRef == nil {
		t.Fail()
	}
}

func TestGetUserActorRef(t *testing.T) {
	t.Log("starting getUserActor test")

	userRef := getUserActorRef("newUser")

	if userRef == nil {
		t.Fail()
	}
}

func TestAddInterestHandler(t *testing.T) {
	t.Log("starting addInterestHandler test")

	testUser := getTestUserActor()
	_ = gosiris.ActorSystem().RegisterActor(testUserName, testUser, nil)

	actorRef, err := gosiris.ActorSystem().ActorOf(testUserName)
	if err != nil {
		t.FailNow()
	}
	err = actorRef.Tell(gosiris.EmptyContext, addInterestMsg, UserInterest{
		UserId: testUserName,
		Interest: getTestInterest(),
	}, actorRef)
	if err != nil {
		t.FailNow()
	}
	time.Sleep(1500 * time.Millisecond)

	if len(testUser.user.Interests) == 0 {
		t.Log("The interest should have been added to the user")
		t.Fail()
	}
}

const testUserName = "testUser"

func getTestUserActor() *userActor {
	user := User{
		Id:            testUserName,
		Interests:     make([]Interest, 0),
		Trades:        make([]TradeOpportunity, 0),
		Notifications: make([]string, 0),
	}
	actor := userActor{
		Actor: gosiris.Actor{},
		user:  user,
	}
	actor.addHandlers()
	return &actor
}

func getTestInterest() Interest {
	return Interest{
		ItemUserid: "otherUser",
		ItemId: "itemId",
	}
}