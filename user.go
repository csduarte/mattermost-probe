package main

import (
	"encoding/json"
	"io"
)

type User struct {
	ID string
}

func UserFromJSON(data io.Reader) *User {
	decoder := json.NewDecoder(data)
	var u User
	err := decoder.Decode(&u)
	if err == nil {
		return &u
	}
	return nil
}
