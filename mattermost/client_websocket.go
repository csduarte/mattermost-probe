package mattermost

import (
	"time"

	"github.com/mattermost/platform/model"
)

// CreateWS will create the websocket connection
func (c *Client) CreateWS(url string) error {
	wsc, err := NewWSClient(url, c.API.GetAuthToken())
	if err != nil {
		return err
	}
	c.WS = wsc
	return nil
}

// StartWS begin listening, call CreateWS first or suffer a panic
func (c *Client) StartWS() {
	c.WS.Listen()
	go func() {
		for {
			select {
			case resp, ok := <-c.WS.GetEventChannel():
				if !ok {
					c.handleWSError()
				} else {
					c.handleWSEvent(resp)
				}
			}
		}
	}()
}

func (c *Client) handleWSEvent(event *model.WebSocketEvent) {
	// TODO: Stamp Event incoming time, instead of doing it in each probe
	c.SubLock.RLock()
	defer c.SubLock.RUnlock()
	for _, wss := range c.Subs {
		if wss.ShouldNotify(event) {
			wss.Emit(event)
		}
	}
}

// AddSubscription will add the subscription probe to the subs for this client
func (c *Client) AddSubscription(s Subscriber) {
	c.SubLock.Lock()
	defer c.SubLock.Unlock()
	c.Subs = append(c.Subs, s.GetSubscription())
}

func (c *Client) handleWSError() {
	retryCount := 0
	if c.WS.GetListenError() != nil {
		c.LogError("Main WebSocket Error: - %s\n", c.WS.GetListenError().Error())
	} else {
		c.LogError("Main WebSocket Error: - Connection closed from server")
	}
	for {
		if retryCount > 0 {
			c.LogInfo("Sleeping for %d seconds before socket reconnect attempt", retryCount)
			time.Sleep(time.Duration(retryCount) * time.Second)
		}
		c.LogInfo("WebSocket attempting to reconnect")
		if err := c.WS.Connect(); err != nil {
			c.LogInfo("Retry WebSocket Error: - %s\n", err.Error())
			retryCount++
			continue
		}
		c.WS.Listen()
		c.LogInfo("WebSocket Reconnected")
		break
	}
}
