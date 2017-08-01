package config

// SearchConfig includes the various resources to be pinged in probe.
type SearchConfig struct {
	Enabled        bool     `yaml:"enabled"`
	Frequency      float64  `yaml:"frequency_sec"`
	UserEnabled    bool     `yaml:"user_enabled"`
	UserTerms      []string `yaml:"user_terms"`
	UserMinimum    int      `yaml:"user_minimum"`
	ChannelEnabled bool     `yaml:"channel_enabled"`
	ChannelTerms   []string `yaml:"channel_terms"`
	ChannelMinimum int      `yaml:"channel_minimum"`
}
