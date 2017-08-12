package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/probe"
	"github.com/csduarte/mattermost-probe/util"
	"github.com/prometheus/common/log"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	applicationStart()
	select {}
}

func applicationStart() {
	//TODO: Add PID check for multiple process

	flagConfig := config.GetFlags()
	log.Info("Application Started")
	log.Info(flagConfig)

	args := flagConfig.Args
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "version" {
			log.Infof("Version X.X.X")
			os.Exit(0)
		} else {
			log.Errorf("application launched with unrecognized arguments %q", args)
			os.Exit(1)
		}
	}

	var log *logrus.Logger
	log, err := util.NewFileLogger(flagConfig.LogLocation, flagConfig.Verbose)
	if err != nil {
		applicationExit(log, err.Error())
	}

	var mlog *logrus.Logger
	if len(flagConfig.MetricsLocation) > 0 {
		mlog, err = util.NewFileLogger(flagConfig.MetricsLocation, false)
		if err != nil {
			applicationExit(log, err.Error())
		}
	}

	cfg := config.Config{}
	file, err := ioutil.ReadFile(flagConfig.ConfigLocation)
	if err != nil {
		applicationExit(log, "Config read error - %s", err.Error())
	}
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		applicationExit(log, "Failed to load config - %s", err.Error())
	}

	if err := cfg.Validate(); err != nil {
		applicationExit(log, "Config error - %s", err.Error())
	}

	server := metrics.NewServer(log, mlog)
	go server.Listen(cfg.BindAddr, cfg.Port)

	c1 := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel, log)
	if err := c1.Establish(cfg.WSHost, cfg.UserA); err != nil {
		applicationExit(log, "Could not establish client 1 - %s", err.Error())
	}

	c2 := mattermost.NewClient(cfg.Host, cfg.TeamID, server.ReportChannel, log)
	if err := c2.Establish(cfg.WSHost, cfg.UserB); err != nil {
		applicationExit(log, "Could not establish client 2 - %s", err.Error())
	}

	probes := probe.NewProbes(cfg, server.ReportChannel, c1, c2)
	if err := probe.SetupProbes(probes, log); err != nil {
		applicationExit(log, err.Error())
	}
	if err := probe.StartProbes(probes, log); err != nil {
		applicationExit(log, err.Error())
	}
	if len(probes) == 0 {
		log.Warn("No probes enabled.")
	} else {
		log.Info("Application running")
	}
}

func applicationExit(log *logrus.Logger, format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}
