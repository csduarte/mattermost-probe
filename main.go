package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/probe"
	"github.com/csduarte/mattermost-probe/util"
	yaml "gopkg.in/yaml.v2"
)

func main() {

	var osSignal = make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGTERM)
	signal.Notify(osSignal, syscall.SIGINT)

	applicationStart()

	select {
	case sig := <-osSignal:
		log.Printf("Application closing from sig %s", sig)
	}
}

func applicationStart() {

	flagConfig := config.GetFlags()

	args := flagConfig.Args
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "version" {
			fmt.Println(Version)
			os.Exit(0)
		} else {
			log.Printf("application launched with unrecognized arguments %q\n", args)
			os.Exit(1)
		}
	}

	log.Println("Application Started")
	log.Println("Version:", Version)
	log.Println(flagConfig)

	var log *logrus.Logger
	log, err := util.NewFileLogger(flagConfig.LogLocation)
	if err != nil {
		applicationExit(log, err.Error())
	}

	var mlog *logrus.Logger
	if len(flagConfig.MetricsLocation) > 0 {
		mlog, err = util.NewFileLogger(flagConfig.MetricsLocation)
		if err != nil {
			applicationExit(log, err.Error())
		}
	}

	log.Info("Flag Config", flagConfig)

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
		applicationExit(log, "Failed to setup probes - %s", err.Error())
	}
	if err := probe.StartProbes(probes, log); err != nil {
		applicationExit(log, "Failed to start probes - %s", err.Error())
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
