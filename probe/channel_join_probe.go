package probe

import (
	"fmt"
	"time"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/pkg/errors"
)

// ChannelJoinProbe represent a user joining a channel
type ChannelJoinProbe struct {
	Name        string
	Client      *mattermost.Client
	Config      config.ChannelJoinConfig
	StopChannel chan bool
	Active      bool
}

// NewChannelJoinProbe creates a channel joining probe
func NewChannelJoinProbe(config config.ChannelJoinConfig, client *mattermost.Client) *ChannelJoinProbe {
	p := ChannelJoinProbe{
		Name:        "Channel Join Probe",
		Client:      client,
		Config:      config,
		StopChannel: make(chan bool),
		Active:      false,
	}
	return &p
}

// Setup will run once on application starts
func (p *ChannelJoinProbe) Setup() error {
	if len(p.Config.ChannelID) < 1 && len(p.Config.ChannelName) < 1 {
		return fmt.Errorf("Must set either ChannelID or ChannelName for probe")
	}

	if p.Config.Frequency < 0.2 {
		p.Client.LogInfo("Frequency cannot be set below 0.2, setting to default 1 sec")
		p.Config.Frequency = 1
	} else {
		p.Client.LogInfo("Channel Join Frequency: %f seconds", p.Config.Frequency)
	}

	if len(p.Config.ChannelID) < 1 {
		p.Client.LogInfo("No Channel ID set, attempting to fetch by ChannelName: %s", p.Config.ChannelName)
		err := p.getChannelID(p.Config.ChannelName)
		if err != nil {
			return errors.Wrap(err, "could not get channel id")
		}
	}

	return nil
}

// Start will kick off the probe
func (p *ChannelJoinProbe) Start() error {
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
				go p.joinChannel()
			}
		}
	}()

	p.Active = true
	return nil
}

func (p *ChannelJoinProbe) getChannelID(name string) error {
	channel, err := p.Client.GetChannelByName(name)
	if err != nil {
		p.Client.LogError("Probe error", err.Error())
		return err
	}
	p.Config.ChannelID = channel.Id
	return nil
}

func (p *ChannelJoinProbe) joinChannel() {
	err := p.Client.JoinChannel(p.Config.ChannelID)
	if err != nil {
		p.Client.LogError("Channel Join Error:", err.Error())
	}
}

func (p *ChannelJoinProbe) String() string {
	return p.Name
}
