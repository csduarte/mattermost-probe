package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	versionURL = "/api/v3/"
)

const (
	tokenHeader  = "token"
	headerBearer = "BEARER"
	headerAuth   = "Authorization"
)

type apiClient struct {
	HTTPClient *http.Client
	Host       string
	AuthToken  string
	UserID     string
	TeamID     string
	ChannelID  string
}

func NewAPIClient(host string) *apiClient {
	return &apiClient{&http.Client{}, host, "", "", "", ""}
}

func (c *apiClient) doAPIPost(url string, data string) (*http.Response, error) {
	fmt.Println("write address: ", c.Host+versionURL+url)
	req, _ := http.NewRequest("POST", c.Host+versionURL+url, strings.NewReader(data))
	req.Close = true

	if len(c.AuthToken) > 0 {
		req.Header.Set(headerAuth, headerBearer+" "+c.AuthToken)
	}

	if resp, err := c.HTTPClient.Do(req); err != nil {
		return nil, fmt.Errorf(err.Error())
	} else if resp.StatusCode >= 300 {
		defer closeBody(resp)
		return nil, fmt.Errorf("Http Error %v", resp)
	} else {
		return resp, nil
	}
}

func (c *apiClient) CreatePost(p *Post) (*Post, error) {
	createURL := fmt.Sprintf("teams/%v/channels/%v/posts/create", c.TeamID, c.ChannelID)
	r, err := c.doAPIPost(createURL, p.toJSON())
	if err != nil {
		return nil, err
	}
	defer closeBody(r)
	return PostFromJSON(r.Body), nil
}

func (c *apiClient) Login(username, password string) error {
	data := fmt.Sprintf("{\"login_id\": \"%s\", \"password\": \"%s\"}", username, password)
	resp, err := c.doAPIPost("users/login", data)
	if err != nil {
		return err
	}
	defer closeBody(resp)
	c.AuthToken = resp.Header.Get(tokenHeader)
	if u := UserFromJSON(resp.Body); u != nil {
		c.UserID = u.ID
	}
	return nil
}

func (c *apiClient) NewSamplePost() *Post {
	return &Post{
		Message:   "This is a message, really",
		UserID:    c.UserID,
		ChannelID: c.ChannelID,
		CreatedAt: 0,
	}
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}
