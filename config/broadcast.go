package config

// BroadcastConfig holds all the configuration for this channel
type BroadcastConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Frequency   float64 `yaml:"frequency_sec"`
	Cutoff      float64 `yaml:"cutoff_sec"`
	ChannelName string  `yaml:"channelName"`
	ChannelID   string  `yaml:"channelID"`
}
