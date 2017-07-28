package config

// PingConfig includes the various resources to be pinged in probe.
type PingConfig struct {
	Enabled   bool                 `yaml:"enabled"`
	Frequency float64              `yaml:"frequency_sec"`
	Resources []PingResourceConfig `yaml:"resources"`
}

// PingResourceConfig holds the url to ping and if the Mattermost Auth token should be included in header according to v4 API.
type PingResourceConfig struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	IncludeAuth bool   `yaml:"include_auth"`
}
