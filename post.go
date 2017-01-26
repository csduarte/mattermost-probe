package main

import (
	"encoding/json"
	"io"
)

type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	RootID    string `json:"root_id"`
	Message   string `json:"message"`
	CreatedAt int    `json:"create_at"`
}

func PostFromJSON(data io.Reader) *Post {
	decoder := json.NewDecoder(data)
	var p Post
	err := decoder.Decode(&p)
	if err == nil {
		return &p
	}
	return nil
}

func (p *Post) toJSON() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}
