package mattermost

import (
	"github.com/mattermost/platform/model"
)

// WSClient satisfies the APIInterface with the mattermost model api client
type WSClient struct {
	model.WebSocketClient
}

// WSInterface is the required set of mattermost client function
type WSInterface interface {
	Listen()
	GetEventChannel() chan *model.WebSocketEvent
	SetEventChannel(chan *model.WebSocketEvent)
	GetListenError() *model.AppError
	GetResponseChannel() chan *model.WebSocketResponse
	SetResponseChannel(chan *model.WebSocketResponse)
	Connect() *model.AppError
}

// NewWSClient makes a WSInterface suitable websocket client
func NewWSClient(url, token string) (WSInterface, *model.AppError) {
	c, err := model.NewWebSocketClient(url, token)
	if err != nil {
		return nil, err
	}
	wsc := &WSClient{*c}
	return wsc, nil
}

// GetEventChannel will return mattermost ws event channel
func (c *WSClient) GetEventChannel() chan *model.WebSocketEvent {
	return c.EventChannel
}

// SetEventChannel will set mattermost ws event channel
func (c *WSClient) SetEventChannel(ec chan *model.WebSocketEvent) {
	c.EventChannel = ec
}

// GetListenError will return a mattermost listen ws error if any, nil otherwise
func (c *WSClient) GetListenError() *model.AppError {
	return c.ListenError
}

// GetResponseChannel will return mattermost Response channel
func (c *WSClient) GetResponseChannel() chan *model.WebSocketResponse {
	return c.ResponseChannel
}

// SetResponseChannel will set mattermost response channel
func (c *WSClient) SetResponseChannel(rc chan *model.WebSocketResponse) {
	c.ResponseChannel = rc
}
