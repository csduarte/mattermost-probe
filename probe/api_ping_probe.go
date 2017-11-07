package probe

import (
	"time"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
)

// APIPingProbe represent will do
type APIPingProbe struct {
	Name        string
	Client      *mattermost.Client
	Config      config.APIPingConfig
	StopChannel chan bool
	Active      bool
}

// NewAPIPingProbe creates a channel joining probe
func NewAPIPingProbe(config config.APIPingConfig, client *mattermost.Client) *APIPingProbe {
	p := APIPingProbe{
		Name:        "API Ping Probe",
		Client:      client,
		Config:      config,
		StopChannel: make(chan bool),
		Active:      false,
	}
	return &p
}

// Setup will run once on application starts
func (p *APIPingProbe) Setup() error {

	if p.Config.Frequency < 0.2 {
		p.Client.LogInfo("Frequency cannot be set below 0.2, setting to default 1 sec")
		p.Config.Frequency = 1
	} else {
		p.Client.LogInfo("%s Frequency: %v seconds", p.Name, p.Config.Frequency)
	}

	return nil
}

// Start will kick off the probe
func (p *APIPingProbe) Start() error {
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
				go p.PingAPI()
			}
		}
	}()

	p.Active = true
	return nil
}

// PingAPI call client api
func (p *APIPingProbe) PingAPI() {
	err := p.Client.PingAPI()
	if err != nil {
		p.Client.LogError("failed to ping api - %s", err.Error())
	}
}

func (p *APIPingProbe) String() string {
	return p.Name
}
