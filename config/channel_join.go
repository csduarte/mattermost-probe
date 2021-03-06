package config

// ChannelJoinConfig represents the configuration for ChannelJoinProbe
type ChannelJoinConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Frequency   float64 `yaml:"frequency_sec"`
	ChannelName string  `yaml:"channel_name"`
	ChannelID   string  `yaml:"channel_id"`
}
