package config

// LoginProbe will login every X seconds.
type LoginProbe struct {
	Enabled   bool    `yaml:"enabled"`
	Frequency float64 `yaml:"frequency_sec"`
}
