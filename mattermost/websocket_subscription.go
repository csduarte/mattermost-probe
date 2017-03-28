package mattermost

import (
	"strings"

	"github.com/mattermost/platform/model"
)

// WebSocketSubscription represents a subscription to websocket events. Will eventually neeed responses too.
type WebSocketSubscription struct {
	EventChanel     chan *model.WebSocketEvent
	RequireAllMatch bool
	EventTypes      []string
	ChannelIDs      []string
	UserIDs         []string
}

// NewWebsocketSubcription returns the most common subscription, still to be determined
func NewWebsocketSubcription(ec chan *model.WebSocketEvent) *WebSocketSubscription {
	wss := WebSocketSubscription{
		ec,
		true,
		[]string{},
		[]string{},
		[]string{},
	}
	return &wss
}

// Emit will pass the event back to the
func (wss *WebSocketSubscription) Emit(event *model.WebSocketEvent) {
	wss.EventChanel <- event
}

// ShouldNotify returns true if the event matches the subscription in at least one way
func (wss *WebSocketSubscription) ShouldNotify(event *model.WebSocketEvent) bool {
	channelMatch := false

	// fmt.Printf("Channels %v in %v?\n", event.Broadcast.ChannelId, wss.ChannelIDs)
	for _, cID := range wss.ChannelIDs {
		if event.Broadcast != nil && cID == event.Broadcast.ChannelId {
			channelMatch = true
		}
	}

	eventMatch := false
	// fmt.Printf("Event: %v in %v?\n", event.Event, wss.EventTypes)
	for _, t := range wss.EventTypes {
		if t == event.Event {
			eventMatch = true
		}
	}

	// Get userID from event data
	userID := ""
	if event.Event == model.WEBSOCKET_EVENT_POSTED {
		post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
		userID = post.UserId
	} else if ui, ok := event.Data["user_id"].(string); ok {
		userID = ui
	}

	userMatch := false
	// fmt.Printf("User %v in %v?\n", userID, wss.UserIDs)
	for _, uID := range wss.UserIDs {
		if uID == userID {
			userMatch = true
		}
	}

	// fmt.Printf("C: %t, U: %t, E: %t\n", channelMatch, userMatch, eventMatch)
	if wss.RequireAllMatch {
		return channelMatch && userMatch && eventMatch
	}
	return channelMatch || userMatch || eventMatch
}
