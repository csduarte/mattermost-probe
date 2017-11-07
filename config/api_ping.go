package config

// APIPingConfig includes the various resources to be pinged in probe.
type APIPingConfig struct {
	Enabled   bool    `yaml:"enabled"`
	Frequency float64 `yaml:"frequency_sec"`
}
