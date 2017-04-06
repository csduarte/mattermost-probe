package mattermost

import (
	"io"
	"net/http"

	"github.com/mattermost/platform/model"
)

// APIClient satisfies the APIInterface with the mattermost model api client
type APIClient struct {
	model.Client
}

// APIInterface is the required set of mattermost client function
type APIInterface interface {
	GetPing() (map[string]string, *model.AppError)
	Login(string, string) (*model.Result, *model.AppError)
	GetChannelByName(string) (*model.Result, *model.AppError)
	JoinChannel(string) (*model.Result, *model.AppError)
	GetFile(string) (io.ReadCloser, *model.AppError)
	CreatePost(*model.Post) (*model.Result, *model.AppError)
	GetTeamID() string
	SetTeamID(string)
	SetTransport(http.RoundTripper)
	GetTransport() http.RoundTripper
	GetAuthToken() string
}

// NewAPIClient returns a new API Client
func NewAPIClient(url string) APIInterface {
	return &APIClient{*model.NewClient(url)}
}

// SetTransport will set the http client round tripper
func (c *APIClient) SetTransport(rt http.RoundTripper) {
	c.HttpClient.Transport = rt
}

// GetTransport will get the http client round tripper
func (c *APIClient) GetTransport() http.RoundTripper {
	return c.HttpClient.Transport
}

// GetAuthToken returns the mattermost auth token
func (c *APIClient) GetAuthToken() string {
	return c.AuthToken
}

// GetTeamID will fetch mattermost client team id
func (c *APIClient) GetTeamID() string {
	return c.TeamId
}

// SetTeamID will set mattermost client team id
func (c *APIClient) SetTeamID(s string) {
	c.SetTeamId(s)
}
