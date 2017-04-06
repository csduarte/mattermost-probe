package mattermost

import (
	"fmt"
	"testing"

	"github.com/mattermost/platform/model"
)

type MockWSClient struct {
	ec chan *model.WebSocketEvent
}

func (c *MockWSClient) GetEventChannel() chan *model.WebSocketEvent {
	return nil
}

func (c *MockWSClient) SetEventChannel(ec chan *model.WebSocketEvent) {
	c.ec = ec
}

func (c *MockWSClient) GetListenError() *model.AppError {
	return nil
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

func TestStartWS(t *testing.T) {
	c := Client{}
	c.API = &MockAPIClient{}
	c.WS = &MockWSClient{}
	// starts and listen, but main routine dies at test end
	c.StartWS()
}

func TestHandleWSEvent(t *testing.T) {
	c := Client{}
	c.API = &MockAPIClient{}
	c.WS = &MockWSClient{}
	ec := make(chan *model.WebSocketEvent)
	c.WS.SetEventChannel(ec)

	c.StartWS()

	go func() {
		ec <- &model.WebSocketEvent{}
	}()

	fmt.Println(<-ec)
	// c.WS = mws
	// // starts and listen, but main routine dies at test end
}
