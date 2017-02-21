package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/probe"
	yaml "gopkg.in/yaml.v2"
)

var logger *Logger
var ezLog int

func init() {

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely-typed key-value pairs.
		"url", "url",
		"attempt", "retryNum",
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "https://url")
}

func main() {
	configLocation := flag.String("config", "./config.yaml", "Config location")
	flag.Parse()

	file, err := ioutil.ReadFile(*configLocation)
	if err != nil {
		applicationExit("Config error - " + err.Error())
	}
	cfg := config.Config{}
	yaml.Unmarshal(file, &cfg)

	if err := cfg.Validate(); err != nil {
		applicationExit("Config error - " + err.Error())
	}

	server := metrics.NewServer()
	go server.Listen()

	userA := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel)
	userB := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel)

	// Need real urls
	if err := userA.Establish(cfg.WSHost, cfg.UserA); err != nil {
		applicationExit("Could not establish user A - " + err.Error())
	}

	if err := userB.Establish(cfg.WSHost, cfg.UserB); err != nil {
		applicationExit("Could not establish user B - " + err.Error())
	}

	probes := []probe.Probe{}

	if cfg.BroadcastProbe.Enabled {
		bp := probe.NewBroadcastProbe(&cfg.BroadcastProbe, userA, userB)
		bp.TimingChannel = server.ReportChannel
		probes = append(probes, bp)
	}

	for _, p := range probes {
		if err := p.Setup(); err != nil {
			// TODO: Need a probe get name function
			applicationExit("Could not setup probe - " + err.Error())
		}
	}

	for _, p := range probes {
		if err := p.Start(); err != nil {
			// TODO: Need a probe get name function
			applicationExit("Could not start probe - " + err.Error())
		}
	}

	fmt.Println("Inital Setup Complete")
	select {}
}

func applicationExit(msg string) {
	fmt.Println("Application Error - ", msg)
	os.Exit(1)
}
