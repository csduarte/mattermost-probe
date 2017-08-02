package mattermost

import (
	"net/http"

	"github.com/gorilla/websocket"
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
	c, err := newMattermostWebsocketClient(url, token)
	if err != nil {
		return nil, err
	}
	wsc := &WSClient{*c}
	return wsc, nil
}

func newMattermostWebsocketClient(url, authToken string) (*model.WebSocketClient, *model.AppError) {
	h := http.Header{
		"Origin": []string{url},
	}
	conn, _, err := websocket.DefaultDialer.Dial(url+model.API_URL_SUFFIX+"/users/websocket", h)
	if err != nil {
		return nil, model.NewLocAppError("NewWebSocketClient", "model.websocket_client.connect_fail.app_error", nil, err.Error())
	}

	client := &model.WebSocketClient{
		Url:             url,
		ApiUrl:          url + model.API_URL_SUFFIX,
		Conn:            conn,
		AuthToken:       authToken,
		Sequence:        1,
		EventChannel:    make(chan *model.WebSocketEvent, 100),
		ResponseChannel: make(chan *model.WebSocketResponse, 100),
		ListenError:     nil,
	}

	client.SendMessage(model.WEBSOCKET_AUTHENTICATION_CHALLENGE, map[string]interface{}{"token": authToken})

	return client, nil
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
