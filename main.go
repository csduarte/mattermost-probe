package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/probe"
	"github.com/csduarte/mattermost-probe/util"
	yaml "gopkg.in/yaml.v2"
)

var log *logrus.Logger

func main() {

	//TODO: Add PID check for multiple process

	var configLocation, logLocation, metricLocation string
	var verbose bool
	flag.StringVar(&configLocation, "config", "./config.yaml", "Config location")
	flag.StringVar(&logLocation, "log", "./mattermost-probe.log", "Log Location, default")
	flag.StringVar(&metricLocation, "metrics", "", "Location for metric logs")
	flag.BoolVar(&verbose, "verbose", false, "Set Log level to debug")
	flag.Parse()

	// TODO: move log establish into config
	log = util.NewFileLogger(logLocation, verbose)

	log.Infof("Application Started")
	log.Infof("Config Location: %s", configLocation)
	log.Infof("Log Location: %s", logLocation)
	log.Infof("Metric Location: %s", metricLocation)

	file, err := ioutil.ReadFile(configLocation)
	if err != nil {
		applicationExit("Config error - " + err.Error())
	}
	cfg := config.Config{}
	yaml.Unmarshal(file, &cfg)

	if err := cfg.Validate(); err != nil {
		applicationExit("Config error - " + err.Error())
	}

	// TODO: Move to server startup metrics package
	server := metrics.NewServer(log, metricLocation)
	go server.Listen(cfg.BindAddr, cfg.Port)

	userA := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel, log)
	userB := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel, log)

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

	if cfg.ChannelJoinProbe.Enabled {
		cjp := probe.NewChannelJoinProbe(&cfg.ChannelJoinProbe, userA)
		probes = append(probes, cjp)
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

	log.Info("All probes established & started")
	select {}
}

func applicationExit(msg string) {
	log.Errorf(msg)
	os.Exit(1)
}
