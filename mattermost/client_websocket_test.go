package mattermost

import (
	"testing"
	"time"

	"github.com/mattermost/platform/model"
)

type MockWSClient struct {
	ec chan *model.WebSocketEvent
}

func (c *MockWSClient) GetEventChannel() chan *model.WebSocketEvent {
	return c.ec
}

func (c *MockWSClient) SetEventChannel(ec chan *model.WebSocketEvent) {
	c.ec = ec
}

func (c *MockWSClient) GetListenError() *model.AppError {
	return model.NewLocAppError("", "", nil, "")
}

func (c *MockWSClient) GetResponseChannel() chan *model.WebSocketResponse {
	return nil
}

func (c *MockWSClient) SetResponseChannel(rc chan *model.WebSocketResponse) {

}

func (c *MockWSClient) Connect() *model.AppError {
	return nil
}

func (c *MockWSClient) Listen() {

}

func TestHandleWSEvent(t *testing.T) {

	success := make(chan bool)
	cutoff := time.NewTimer(100 * time.Microsecond)

	c := Client{}
	wsEvents := make(chan *model.WebSocketEvent)
	probeEvents := make(chan *model.WebSocketEvent)

	p := &MockSubscriptionProbe{EventChannel: probeEvents}
	p.TargetEvent = model.WEBSOCKET_EVENT_POSTED

	c.API = &MockAPIClient{}
	c.WS = &MockWSClient{}
	c.WS.SetEventChannel(wsEvents)

	c.AddSubscription(p)
	c.StartWS()

	// Act as temp probe and listen
	go func() {
		if e := <-p.EventChannel; e != nil {
			success <- true
		}
		close(success)
	}()

	// Act as Mattermost and send event
	go func() {
		wsEvents <- &model.WebSocketEvent{Event: model.WEBSOCKET_EVENT_POSTED}
		// intentionally not closed to avoid err handling
	}()

	select {
	case <-success:
		break
	case <-cutoff.C:
		t.Fatal("No event before cutoff from subscription")
	}
}

// UNUSED TEST
// func TestHandleWSError(t *testing.T) {

// 	success := make(chan bool)

// 	probeEvents := make(chan *model.WebSocketEvent)
// 	p := &MockSubscriptionProbe{
// 		EventChannel: probeEvents,
// 		TargetEvent:  model.WEBSOCKET_EVENT_POSTED,
// 	}

// 	wsEvents := make(chan *model.WebSocketEvent)
// 	c := Client{
// 		API: &MockAPIClient{},
// 		WS:  &MockWSClient{},
// 	}
// 	c.WS.SetEventChannel(wsEvents)
// 	c.AddSubscription(p)
// 	c.StartWS()

// 	// Act as probe and listen
// 	go func() {
// 		if e := <-p.EventChannel; e != nil {
// 			success <- true
// 		}
// 		close(success)
// 	}()

// 	// Act as Mattermost and close ws
// 	go func() {
// 		close(wsEvents)
// 	}()
// }
