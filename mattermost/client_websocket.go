package mattermost

import "github.com/mattermost/platform/model"

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
			case resp := <-c.WS.EventChannel:
				c.handleWSEvent(resp)
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
