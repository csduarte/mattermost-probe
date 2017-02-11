package config

import "time"

// BroadcastConfig holds all the configuration for this channel
type BroadcastConfig struct {
	Enabled     bool          `yaml:"enabled"`
	Frequency   time.Duration `yaml:"frequency_ms"`
	ChannelName string        `yaml:"channelName"`
	ChannelID   string        `yaml:"channelID"`
}
