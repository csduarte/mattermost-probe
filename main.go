package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/probe"
	"github.com/csduarte/mattermost-probe/util"
	yaml "gopkg.in/yaml.v2"
)

var log *zap.SugaredLogger
var output *zap.SugaredLogger

func main() {

	configLocation := *flag.String("config", "./config.yaml", "Config location")
	logLocation := *flag.String("log", "./mattermost-probe.log", "Log Location, default")
	outputLocation := *flag.String("output", "", "Location for metric logs")
	flag.Parse()

	log, err := util.NewEasyLogger(logLocation)
	if err != nil {
		applicationExit(fmt.Sprintf("Failed to create logger - %s", err.Error()))
	}

	log.Infof("Application Started")
	log.Infof("Config Location: %s", configLocation)
	log.Infof("Log Location: %s", logLocation)
	log.Infof("Ouptut Locaitn: %s", outputLocation)

	if len(outputLocation) > 0 {
		// TODO: init metrics logger
	}

	file, err := ioutil.ReadFile(configLocation)
	if err != nil {
		applicationExit("Config error - " + err.Error())
	}
	cfg := config.Config{}
	yaml.Unmarshal(file, &cfg)

	if err := cfg.Validate(); err != nil {
		applicationExit("Config error - " + err.Error())
	}

	server := metrics.NewServer(log, outputLocation)
	go server.Listen(cfg.BindAddr, cfg.Port)

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

	log.Info("Startup Complete")
	select {}
}

func applicationExit(msg string) {
	log.Errorf(msg)
	os.Exit(1)
}
