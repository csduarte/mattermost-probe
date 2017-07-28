package probe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
)

// PingProbe represent will do
type PingProbe struct {
	Name        string
	Client      *mattermost.Client
	Config      config.PingConfig
	StopChannel chan bool
	Active      bool
}

// NewPingProbe creates a channel joining probe
func NewPingProbe(config config.PingConfig, client *mattermost.Client) *PingProbe {
	p := PingProbe{
		Name:        "Ping Probe",
		Client:      client,
		Config:      config,
		StopChannel: make(chan bool),
		Active:      false,
	}
	return &p
}

// Setup will run once on application starts
func (p *PingProbe) Setup() error {

	if p.Config.Frequency < 0.2 {
		p.Client.LogInfo("Frequency cannot be set below 0.2, setting to default 1 sec")
		p.Config.Frequency = 1
	} else {
		p.Client.LogInfo("Ping Frequency: %v seconds", p.Config.Frequency)
	}

	return nil
}

// Start will kick off the probe
func (p *PingProbe) Start() error {
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
				go p.CheckResources()
			}
		}
	}()

	p.Active = true
	return nil
}

// CheckResources will send a ping to each configured resource and ensure valid response
func (p *PingProbe) CheckResources() {
	bearer := fmt.Sprintf("Bearer %s", p.Client.API.GetAuthToken())
	hc := p.Client.API.GetHTTPClient()
	for _, r := range p.Config.Resources {
		req, err := http.NewRequest("GET", r.URL, nil)
		if err != nil {
			p.Client.LogError("Failed to form ping request to", r.URL)
		}
		if r.IncludeAuth {
			req.Header.Add("Authorization", bearer)
		}
		req.Header.Add("PingProbe", r.Name)

		// res, error caught on trainsport layer and counted
		hc.Do(req)
	}
}

func (p *PingProbe) String() string {
	return p.Name
}
