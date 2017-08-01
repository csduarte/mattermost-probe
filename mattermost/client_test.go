package mattermost

import (
	"io"
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/mattermost/platform/model"
)

type MockAPIClient struct {
	shouldFail        bool
	shouldFailPing    bool
	shouldFailLogin   bool
	shouldFailStartWS bool
}

func (c *MockAPIClient) GetPing() (map[string]string, *model.AppError) {
	if c.shouldFail || c.shouldFailPing {
		return map[string]string{}, model.NewLocAppError("", "", nil, "")
	}
	return map[string]string{}, nil
}

func (c *MockAPIClient) Login(string, string) (*model.Result, *model.AppError) {
	if c.shouldFail || c.shouldFailLogin {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return &model.Result{Data: &model.User{}}, nil
}

func (c *MockAPIClient) GetChannelByName(string) (*model.Result, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return &model.Result{Data: &model.Channel{}}, nil
}

func (c *MockAPIClient) JoinChannel(string) (*model.Result, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return &model.Result{Data: &model.Channel{}}, nil
}

func (c *MockAPIClient) SearchUsers(params model.UserSearch) (*model.Result, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return &model.Result{Data: model.User{}}, nil
}

func (c *MockAPIClient) SearchMoreChannels(channelSearch model.ChannelSearch) (*model.Result, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return &model.Result{Data: model.ChannelList{}}, nil
}

func (c *MockAPIClient) GetFile(string) (io.ReadCloser, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return nil, nil
}

func (c *MockAPIClient) CreatePost(*model.Post) (*model.Result, *model.AppError) {
	if c.shouldFail {
		return nil, model.NewLocAppError("", "", nil, "")
	}
	return nil, nil
}

func (c *MockAPIClient) SetTeamID(string) {
}

func (c *MockAPIClient) GetTeamID() string {
	return ""
}

func (c *MockAPIClient) SetTransport(http.RoundTripper) {
}

func (c *MockAPIClient) GetTransport() http.RoundTripper {
	return nil
}
func (c *MockAPIClient) GetAuthToken() string {
	return ""
}

func (c *MockAPIClient) GetHTTPClient() *http.Client {
	return nil
}

func TestNewClient(t *testing.T) {
	var tc chan metrics.Report
	var log logrus.Logger
	mockID := "teamID"
	c := NewClient("", mockID, tc, &log)
	if c.API.GetTeamID() != mockID {
		t.Fatalf("Failed to set teamID expected: %v got: %v", c.API.GetTeamID(), mockID)
	}
	if _, ok := c.API.GetTransport().(*metrics.TimedRoundTripper); ok {
		t.Fatal("TimedRoundTripper should *not* be set without timing channel")
	}

	tc = make(chan metrics.Report)
	c = NewClient("", "teamID", tc, &log)
	if _, ok := c.API.GetTransport().(*metrics.TimedRoundTripper); !ok {
		t.Fatal("TimedRoundTripper should be set")
	}
}

func TestClientEstablish(t *testing.T) {
	// c := Client{}
	// mc := &MockAPIClient{}
	// c.API = mc
	// if err := c.Establish("", config.Credentials{}); err != nil {
	// 	t.Fatal("Establish should not return an error")
	// }

}

func TestPingAPI(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	if err := c.PingAPI(); err != nil {
		t.Fatal("PingAPI should return no error")
	}

	mc.shouldFail = true
	if err := c.PingAPI(); err == nil {
		t.Fatal("PingAPI should return error when fails")
	}
}

func TestLogin(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	if err := c.Login(config.Credentials{}); err != nil {
		t.Fatal("Login should *not* return an error")
	}
	if c.User == nil {
		t.Fatal("User should be set after login")
	}
	mc.shouldFail = true
	if err := c.Login(config.Credentials{}); err == nil {
		t.Fatal("Login should return an error")
	}
}

func TestGetChannelByName(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	cl, err := c.GetChannelByName("")
	if err != nil {
		t.Fatal("GetChannelByName should return no error")
	}
	if cl == nil {
		t.Fatal("GetChannelByName should return a channel")
	}

	mc.shouldFail = true
	_, err = c.GetChannelByName("")
	if err == nil {
		t.Fatal("GetChannelByName should return an error")
	}
}

func TestJoinChannel(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	err := c.JoinChannel("")
	if err != nil {
		t.Fatal("JoinChannel should return no error")
	}

	mc.shouldFail = true
	err = c.JoinChannel("")
	if err == nil {
		t.Fatal("JoinChannel should return an error")
	}
}

func TestGetFile(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	err := c.GetFile("")
	if err != nil {
		t.Fatal("GetFile should return no error")
	}

	mc.shouldFail = true
	err = c.GetFile("")
	if err == nil {
		t.Fatal("GetFile should return an error")
	}
}

func TestCreatePost(t *testing.T) {
	c := Client{}
	mc := &MockAPIClient{}
	c.API = mc
	err := c.CreatePost(&model.Post{})
	if err != nil {
		t.Fatal("CreatePost should return no error")
	}

	mc.shouldFail = true
	err = c.CreatePost(&model.Post{})
	if err == nil {
		t.Fatal("CreatePost should return an error")
	}
}

func TestLogInfoAndError(t *testing.T) {
	c := Client{}
	c.LogError("Test Error Log")
	c.LogInfo("Test Info Log")
	c.Log = &logrus.Logger{}
	c.LogError("Test Error Log")
	c.LogInfo("Test Info Log")
}
