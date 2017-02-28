package mattermost

import (
	"time"

	"github.com/mattermost/platform/model"
)

// StartWS will create the websocket connection and begin listening
func (c *Client) StartWS(url string) error {
	wsClient, err := model.NewWebSocketClient(url, c.API.AuthToken)
	if err != nil {
		return err
	}
	c.WS = wsClient
	c.WS.Listen()

	go func() {
		for {
			select {
			case resp, ok := <-c.WS.EventChannel:
				if !ok {
					c.handleWSError()
				} else {
					c.handleWSEvent(resp)
				}
			}
		}
	}()
	return nil
}

func (c *Client) handleWSEvent(event *model.WebSocketEvent) {
	for _, wss := range c.Subs {
		if wss.ShouldNotify(event) {
			wss.Emit(event)
		}
	}
}

// AddSubscription will add the subscription probe to the subs for this client
func (c *Client) AddSubscription(ps SubscriptionProbe) {
	c.Subs = append(c.Subs, ps.GetSubscription())
}

func (c *Client) handleWSError() {
	retryCount := 0
	if c.WS.ListenError != nil {
		c.LogError("Main WebSocket Error: - %v\n", c.WS.ListenError.Error())
	} else {
		c.LogError("Main WebSocket Error: - Connection closed from server")
	}
	for {
		if retryCount > 0 {
			time.Sleep(time.Duration(retryCount) * time.Second)
		}
		c.WS.EventChannel = make(chan *model.WebSocketEvent, 100)
		c.WS.ResponseChannel = make(chan *model.WebSocketResponse, 100)
		c.LogError("WebSocket attempting to reconnected")
		if err := c.WS.Connect(); err != nil {
			retryCount++
			continue
		}
		c.WS.Listen()
		c.LogError("WebSocket Reconnected")
		break
	}
}
