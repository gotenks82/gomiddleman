package models

import "github.com/jinzhu/copier"

type (
	Interest struct {
		ItemUserid       	string `json:"itemUserid"`
		ItemId     		 	string `json:"itemId"`
		ItemName  			string `json:"itemName"`
		ItemImageUrl      	string `json:"itemImageUrl"`
		ItemUrl 			string `json:"itemUrl"`
		ItemPrice        	string `json:"itemPrice"`
	}

	User struct {
		Id string `json:"id"`
		Interests []Interest `json:"interests"`
		Trades []TradeOpportunity `json:"tradeOpportunities"`
		Notifications []string `json:"notifications"`
	}

	UserInterest struct {
		UserId   string `json:"userId"`
		Interest Interest `json:"interests"`
	}

	TradeOpportunity struct {
		Id string `json:"id"`
		RootUserId string `json:"rootUserId"`
		Steps []UserInterest `json:"steps"`
	}

	NotificationsList struct {
		Notifications []string `json:"notifications"`
	}
)

func buildInterest() Interest {
	return Interest{
		ItemName: "",
		ItemImageUrl: "",
		ItemUrl: "",
		ItemPrice: "",
	}
}

func (user *User) AddInterest(interest Interest) {
	user.Interests = append(user.Interests, interest)
}

func (user *User) AddTrade(trade TradeOpportunity) {
	user.Trades = append(user.Trades, trade)
}

func (user *User) AddNotification(notification string) {
	user.Notifications = append(user.Notifications, notification)
}

func (user *User) GetAndResetNotifications() []string {
	notifications := user.Notifications
	user.Notifications = make([]string, 0)
	return notifications
}

func (src TradeOpportunity) clone() TradeOpportunity {
	dest := TradeOpportunity{}
	_ = copier.Copy(&dest, &src)
	return dest
}

func (trade TradeOpportunity) addStep(interest UserInterest) TradeOpportunity {
	trade.Steps = append(trade.Steps, interest)
	return trade
}

func (trade TradeOpportunity) CopyWithStep(interest UserInterest) TradeOpportunity {
	return trade.clone().addStep(interest)
}

func (trade TradeOpportunity) IsUserInvolved(userId string) bool {
	for _, userInterest := range trade.Steps {
		if userInterest.UserId == userId {
			return true
		}
	}
	return false
}

func (trade TradeOpportunity) WasReceivedBy(userId string) bool {
	for _, userInterest := range trade.Steps {
		if userInterest.Interest.ItemUserid == userId {
			return true
		}
	}
	return false
}

func (trade TradeOpportunity) IsComplete() bool {
	return trade.WasReceivedBy(trade.RootUserId)
}