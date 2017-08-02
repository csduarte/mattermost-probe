package mattermost

import (
	"fmt"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/mattermost/platform/model"
)

// Client structure holds mattermost api, websocket, and user objects
type Client struct {
	API  APIInterface
	WS   WSInterface
	User *model.User
	Subs []*WebSocketSubscription
	Log  *logrus.Logger
}

// NewClient generateds a new API and WebSocket Client}
func NewClient(url, teamID string, tc chan metrics.Report, log *logrus.Logger) *Client {
	c := Client{
		NewAPIClient(url),
		nil,
		nil,
		[]*WebSocketSubscription{},
		log,
	}
	c.API.SetTeamID(teamID)
	if tc != nil {
		c.API.SetTransport(metrics.NewTimedRoundTripper(tc, log))
	}
	return &c
}

// Establish will ping the server, login, and create the websocket connection
func (c *Client) Establish(socketURL string, creds config.Credentials) error {
	if _, err := c.API.GetPing(); err != nil {
		return fmt.Errorf("Server Down: %v", err.Error())
	}

	if err := c.Login(creds); err != nil {
		return fmt.Errorf("Failed to login: %v", err.Error())
	}

	if err := c.CreateWS(socketURL); err != nil {
		return fmt.Errorf("Failed to connect ws: %v", err.Error())
	}
	c.StartWS()
	return nil
}

// Listen will prompt the websocket to start a listening routine
func (c *Client) Listen() {
	c.StartWS()
}

// PingAPI will call the ping endpoint
func (c *Client) PingAPI() error {
	if _, err := c.API.GetPing(); err != nil {
		return err
	}
	return nil
}

// Login will the login endpoint
func (c *Client) Login(creds config.Credentials) error {
	results, err := c.API.Login(creds.Email, creds.Password)
	if err != nil {
		return err
	}
	data, ok := results.Data.(*model.User)
	if !ok {
		return fmt.Errorf("Login could not extract channel user data from response")
	}
	c.User = data
	return nil
}

// GetChannelByName will get a channel by the system name, not display name
func (c *Client) GetChannelByName(name string) (*model.Channel, error) {
	results, err := c.API.GetChannelByName(name)
	if err != nil {
		return nil, err
	}
	data, ok := results.Data.(*model.Channel)
	if !ok {
		return nil, fmt.Errorf("GetChannelByName could not extract channel data from response")
	}
	return data, nil
}

// JoinChannel joines the client's user to a channel by channelID
func (c *Client) JoinChannel(id string) error {
	if _, err := c.API.JoinChannel(id); err != nil {
		return err
	}
	return nil
}

// SearchChannels takes a term and returns a slice of MM channels.
func (c *Client) SearchChannels(term string) (*model.ChannelList, error) {
	cs := model.ChannelSearch{Term: term}
	results, err := c.API.SearchMoreChannels(cs)
	if err != nil {
		return nil, err
	}
	cl, ok := results.Data.(*model.ChannelList)
	if !ok {
		return nil, fmt.Errorf("Client SearchChannel could not extract data from response")
	}
	val := reflect.ValueOf(cl)
	if val.IsNil() {
		return nil, fmt.Errorf("Client SearchChannels bad response caused nil channel list")
	}
	return cl, nil
}

// SearchUsers takes a term string and eturns a slice of MM Users.
func (c *Client) SearchUsers(term string) ([]*model.User, error) {
	us := model.UserSearch{Term: term}
	results, err := c.API.SearchUsers(us)
	if err != nil {
		return nil, err
	}
	ul, ok := results.Data.([]*model.User)
	if !ok {
		return nil, fmt.Errorf("Client SearchUsers could not extract data from response")
	}
	return ul, nil
}

// GetFile will fetch a file by file ID
func (c *Client) GetFile(id string) error {
	if _, err := c.API.GetFile(id); err != nil {
		return err
	}
	return nil
}

// CreatePost sends post to api
func (c *Client) CreatePost(post *model.Post) error {
	if _, err := c.API.CreatePost(post); err != nil {
		return err
	}
	return nil
}

// LogError is helper for error logging
func (c *Client) LogError(format string, args ...interface{}) {
	if c.Log != nil {
		c.Log.Errorf(format, args)
	}
}

// LogInfo is a helper for info logging
func (c *Client) LogInfo(format string, args ...interface{}) {
	if c.Log != nil {
		c.Log.Infof(format, args)
	}
}
