package probe

import (
	"time"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
)

// LoginProbe represent will do
type LoginProbe struct {
	Name        string
	Client      *mattermost.Client
	Credentials config.Credentials
	Config      config.LoginProbeConfig
	StopChannel chan bool
	Active      bool
}

// NewLoginProbe creates a channel joining probe
func NewLoginProbe(config config.LoginProbeConfig, client *mattermost.Client, creds config.Credentials) *LoginProbe {
	p := LoginProbe{
		Name:        "Login Probe",
		Client:      client,
		Credentials: creds,
		Config:      config,
		StopChannel: make(chan bool),
		Active:      false,
	}
	return &p
}

// Setup will run once on application starts
func (p *LoginProbe) Setup() error {

	if p.Config.Frequency < 1 {
		p.Client.LogInfo("Frequency cannot be set below 1, setting to default 5 sec")
		p.Config.Frequency = 5
	} else {
		p.Client.LogInfo("%s Frequency: %v seconds", p.Name, p.Config.Frequency)
	}

	return nil
}

// Start will kick off the probe
func (p *LoginProbe) Start() error {
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
				go p.Login()
			}
		}
	}()

	p.Active = true
	return nil
}

// Login will fire off an attempt to auth against server, errors are caught by transport layer
func (p *LoginProbe) Login() {
	p.Client.Login(p.Credentials)
	p.Client.Logout()
}

func (p *LoginProbe) String() string {
	return p.Name
}
