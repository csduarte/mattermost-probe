package config

import "fmt"

// Config represents the application config
type Config struct {
	Host             string            `yaml:"host"`
	Port             int               `yaml:"port"`
	TeamID           string            `yaml:"team_id"`
	WSHost           string            `yaml:"ws_host"`
	BindAddr         string            `yaml:"bind_address"`
	UserA            Credentials       `yaml:"user_a"`
	UserB            Credentials       `yaml:"user_b"`
	PingProbe        PingConfig        `yaml:"ping_probe"`
	BroadcastProbe   BroadcastConfig   `yaml:"broadcast_probe"`
	ChannelJoinProbe ChannelJoinConfig `yaml:"channel_join_probe"`
	SearchProbe      SearchConfig      `yaml:"search_probe"`
}

// Validate will ensure that all required values are or default values are set
func (c *Config) Validate() error {
	if len(c.TeamID) < 1 {
		return fmt.Errorf("Must set a 'team_id' for probes to test")
	}
	if len(c.Host) < 1 {
		return fmt.Errorf("Must set a 'host' for probes to login")
	}
	if len(c.WSHost) < 1 {
		return fmt.Errorf("Must set a 'ws_host' for probes to listen to messages")
	}

	if !c.UserA.Valid() {
		return fmt.Errorf("'user_a' missing either email or password")
	}

	if !c.UserB.Valid() {
		return fmt.Errorf("'user_b' missing either email or password")
	}

	return nil
}
