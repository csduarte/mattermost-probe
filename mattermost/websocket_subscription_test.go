package mattermost

import (
	"testing"
	"time"

	"github.com/mattermost/platform/model"
)

func TestEmit(t *testing.T) {
	var cutoff time.Duration = 1
	ech := make(chan *model.WebSocketEvent)
	defer close(ech)

	s := NewWebsocketSubcription(ech)
	tmr := time.NewTimer(cutoff * time.Millisecond)
	go func() {
		s.Emit(&model.WebSocketEvent{})
	}()

	select {
	case <-ech:
		break
	case <-tmr.C:
		t.Fatalf("Emit did complete in complete %d ms", cutoff)
	}
}

func TestShouldNotify(t *testing.T) {
	ech := make(chan *model.WebSocketEvent)
	defer close(ech)

	// No matches
	s := NewWebsocketSubcription(ech)
	e := model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	if s.ShouldNotify(e) != false {
		t.Fatal("ShouldNotify should return false when there are no matches")
	}

	// Test Channel Match (Require any)
	s = NewWebsocketSubcription(ech)
	s.ChannelIDs = append(s.ChannelIDs, "channel")
	s.RequireAllMatch = false
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	if s.ShouldNotify(e) != true {
		t.Fatalf("ShouldNotify should return true when channel matches with require any. expected: %q, actual: %q", s.ChannelIDs, e.Event)
	}

	// Test Event Type Match (Require All)
	s = NewWebsocketSubcription(ech)
	s.EventTypes = append(s.EventTypes, "event")
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	if s.ShouldNotify(e) != false {
		t.Fatalf("ShouldNotify should return false when event type matches with require all. expected: %q, actual: %q", s.EventTypes, e.Event)
	}

	// Test Event Type Match (Require any)
	s = NewWebsocketSubcription(ech)
	s.EventTypes = append(s.EventTypes, "event")
	s.RequireAllMatch = false
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	if s.ShouldNotify(e) != true {
		t.Fatalf("ShouldNotify should return true when event type matches with require any. expected: %q, actual: %q", s.EventTypes, e.Event)
	}

	// Test UserID from Post
	s = NewWebsocketSubcription(ech)
	s.UserIDs = append(s.EventTypes, "user")
	s.RequireAllMatch = false
	p := model.Post{UserId: "user"}
	pJSON := p.ToJson()
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	e.Event = model.WEBSOCKET_EVENT_POSTED
	e.Data = map[string]interface{}{"post": pJSON}
	if s.ShouldNotify(e) != true {
		t.Fatalf("ShouldNotify should return true when post is from matching userid")
	}

	// Test UserID from None Post
	s = NewWebsocketSubcription(ech)
	s.UserIDs = append(s.EventTypes, "user")
	s.RequireAllMatch = false
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	e.Data = map[string]interface{}{"user_id": "user"}
	if s.ShouldNotify(e) != true {
		t.Fatalf("ShouldNotify should return true when data is from matching userid")
	}

	// Test All Match
	// Test UserID from None Post
	s = NewWebsocketSubcription(ech)
	s.UserIDs = append(s.EventTypes, "user")
	s.EventTypes = append(s.EventTypes, "event")
	s.ChannelIDs = append(s.ChannelIDs, "channel")
	e = model.NewWebSocketEvent("event", "team", "channel", "user", map[string]bool{})
	e.Data = map[string]interface{}{"user_id": "user"}
	if s.ShouldNotify(e) != true {
		t.Fatalf("ShouldNotify should return true when data is from matching userid")
	}
}
