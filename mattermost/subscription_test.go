package mattermost 

import "github.com/mattermost/platform/model"

type MockSubscriptionProbe struct {
	EventChannel  chan *model.WebSocketEvent
	TargetUser    string
	TargetChannel string
	TargetEvent   string
}

func (p *MockSubscriptionProbe) GetSubscription() *WebSocketSubscription {
	wss := NewWebsocketSubcription(p.EventChannel)
	wss.RequireAllMatch = false
	wss.UserIDs = append(wss.UserIDs, p.TargetUser)
	wss.ChannelIDs = append(wss.ChannelIDs, p.TargetChannel)
	wss.EventTypes = append(wss.EventTypes, p.TargetEvent)
	return wss
}
