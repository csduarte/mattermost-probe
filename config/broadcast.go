package config

// BroadcastConfig holds all the configuration for this channel
type BroadcastConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Frequency   float64 `yaml:"frequency_sec"`
	Cutoff      float64 `yaml:"cutoff_sec"`
	ChannelName string  `yaml:"channel_name"`
	ChannelID   string  `yaml:"channel_id"`
}
