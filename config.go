package main

type Config struct {
	Cutoff    int
	ChannelID string `yaml:"channelID"`
	TeamID    string `yaml:"teamID"`
	Username  string
	Password  string
	Hosts     []string
}
