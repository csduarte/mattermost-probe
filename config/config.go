package config

import "fmt"

// Config represents the application config
type Config struct {
	TeamID         string          `yaml:"team_id"`
	Host           string          `yaml:"host"`
	WSHost         string          `yaml:"ws_host"`
	UserA          Credentials     `yaml:"user_a"`
	UserB          Credentials     `yaml:"user_b"`
	BroadcastProbe BroadcastConfig `yaml:"broadcast_config"`
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
