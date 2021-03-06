package config

import (
	"flag"
	"fmt"
)

type FlagsConfig struct {
	ConfigLocation  string   `yaml:"config"`
	MetricsLocation string   `yaml:"metrics_location"`
	LogLocation     string   `yaml:"log_location"`
	Args            []string `yaml:"args"`
}

func (c FlagsConfig) String() string {
	ml := "<Disabled>"
	if len(c.MetricsLocation) > 0 {
		ml = c.MetricsLocation
	}
	return fmt.Sprintf(`Flag Configuration ---
File Location(Config):  %s
File Location(Log):     %s
File Location(Metric):  %s
---`,
		c.ConfigLocation, c.LogLocation, ml)
}

func GetFlags() FlagsConfig {
	c := FlagsConfig{}
	flag.StringVar(&c.ConfigLocation, "config", "./config.yaml", "Config location including filename")
	flag.StringVar(&c.LogLocation, "log", "./mattermost-probe.log", "Log Location including filename")
	flag.StringVar(&c.MetricsLocation, "metrics", "", "Metric Log Location including filename")
	flag.Parse()
	c.Args = flag.Args()
	return c
}
