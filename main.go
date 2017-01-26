package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cfg := Config{}
	yaml.Unmarshal(file, &cfg)

	var clients []*apiClient
	for _, h := range cfg.Hosts {
		fmt.Printf("Creating client for %v\n", h)
		c := NewAPIClient(h)
		c.ChannelID = cfg.ChannelID
		c.TeamID = cfg.TeamID

		clients = append(clients, c)

		if err = c.Login(cfg.Username, cfg.Password); err != nil {
			fmt.Printf("Failed login: %v\n", err.Error())
			os.Exit(1)
		}

		wc := NewWriteCheck(c)
		wc.Start()
	}

	select {}
}
