package main

import (
	"fmt"
	"os"
	"time"
)

type WriteCheck struct {
	Client      *apiClient
	Active      bool
	StopChannel chan bool
	Count       int
}

func NewWriteCheck(c *apiClient) *WriteCheck {
	return &WriteCheck{c, false, make(chan bool), 0}
}

func (wc *WriteCheck) Start() {
	if wc.Active {
		return
	}
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-wc.StopChannel:
				return
			case <-ticker.C:
				fmt.Println(time.Now())
				wc.SendWrite()
			}
		}
	}()
}

func (wc *WriteCheck) Stop() {
	wc.StopChannel <- true
}

func (wc *WriteCheck) SendWrite() {
	p := wc.Client.NewSamplePost()
	fmt.Println("Creating post")
	_, err := wc.Client.CreatePost(p)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
