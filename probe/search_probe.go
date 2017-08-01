package probe

import (
	"fmt"
	"time"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
)

type searchType int

const (
	searchUsers searchType = iota
	searchChannels
)

// SearchProbe represent will do
type SearchProbe struct {
	Name           string
	Client         *mattermost.Client
	Config         config.SearchConfig
	ReportChannel  chan metrics.Report
	StopChannel    chan bool
	Active         bool
	ChannelTermGen chan string
	UserTermGen    chan string
}

// NewSearchProbe creates a channel joining probe
func NewSearchProbe(config config.SearchConfig, client *mattermost.Client) *SearchProbe {
	p := SearchProbe{
		Name:           "Search Probe",
		Client:         client,
		Config:         config,
		StopChannel:    make(chan bool),
		Active:         false,
		ChannelTermGen: NewGenerator(config.ChannelTerms),
		UserTermGen:    NewGenerator(config.UserTerms),
	}
	return &p
}

// NewGenerator creates a goroutine that will loop through the given array
// and return the next element. The generator will start over after
// returning the entire array.
func NewGenerator(s []string) chan string {
	c := make(chan string)
	l := len(s)
	go func() {
		for {
			if l < 1 {
				c <- ""
			}
			for _, ss := range s {
				c <- ss
			}
		}
	}()
	return c
}

// Setup will run once on application starts
func (p *SearchProbe) Setup() error {

	if p.Config.Frequency < 0.2 {
		p.Client.LogInfo("Frequency cannot be set below 0.2, setting to default 1 sec")
		p.Config.Frequency = 1
	} else {
		p.Client.LogInfo("Ping Frequency: %v seconds", p.Config.Frequency)
	}

	return nil
}

// Start will kick off the probe
func (p *SearchProbe) Start() error {
	if p.Active {
		return nil
	}

	t := time.Duration(p.Config.Frequency * float64(time.Second))
	writeTicker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-p.StopChannel:
				return
			case <-writeTicker.C:
				if p.Config.ChannelEnabled {
					go p.SearchChannels()
				}
				if p.Config.UserEnabled {
					go p.SearchUsers()
				}
			}
		}
	}()

	p.Active = true
	return nil
}

func (p *SearchProbe) SearchUsers() {
	t := <-p.UserTermGen
	min := p.Config.UserMinimum
	// error caught by transport layer
	users, _ := p.Client.SearchUsers(t)
	if len(users) < min {
		p.ReportLowResults(searchUsers, t, len(users))
		return
	}
	return
}

func (p *SearchProbe) SearchChannels() {
	t := <-p.ChannelTermGen
	min := p.Config.ChannelMinimum
	// error caught by transport layer
	channels, _ := p.Client.SearchChannels(t)
	cnt := len(*channels)
	if cnt < min {
		p.ReportLowResults(searchChannels, t, cnt)
	}
	return
}

func (p *SearchProbe) String() string {
	return p.Name
}

func (p *SearchProbe) ReportLowResults(t searchType, term string, c int) {
	if p.ReportChannel == nil {
		return
	}
	switch t {
	case searchUsers:
		p.ReportChannel <- metrics.Report{
			Route:           metrics.RouteUserSearchLowResult,
			DurationSeconds: 0,
			Error:           fmt.Errorf("Search for user %q returned %d result(s). Expected minimum: %d", term, c, p.Config.UserMinimum),
		}
	case searchChannels:
		p.ReportChannel <- metrics.Report{
			Route:           metrics.RouteChannelSearchLowResult,
			DurationSeconds: 0,
			Error:           fmt.Errorf("Search for channel %q returned %d result(s). Expected minimum: %d", term, c, p.Config.ChannelMinimum),
		}
	}
}
