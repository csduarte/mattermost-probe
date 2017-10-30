package probe

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/csduarte/mattermost-probe/util"
	"github.com/mattermost/platform/model"
	uuid "github.com/satori/go.uuid"
)

// BroadcastProbe represents a test where the speaker will broadcast unique messages and the listener will check broadcast time.
type BroadcastProbe struct {
	Name          string
	Speaker       *mattermost.Client
	Listener      *mattermost.Client
	Config        config.BroadcastConfig
	Messages      *util.MessageMap
	EventChannel  chan *model.WebSocketEvent
	ReportChannel chan metrics.Report
	StopChannel   chan bool
	Active        bool
}

// NewBroadcastProbe creates a new base probe
func NewBroadcastProbe(c config.BroadcastConfig, speaker, listener *mattermost.Client) *BroadcastProbe {
	bp := BroadcastProbe{
		"Broadcast Probe",
		speaker,
		listener,
		c,
		util.NewMessageMap(),
		make(chan *model.WebSocketEvent, 10),
		nil,
		make(chan bool),
		false,
	}
	return &bp
}

// Setup will run once on application starts
func (bp *BroadcastProbe) Setup() error {
	if len(bp.Config.ChannelID) < 1 && len(bp.Config.ChannelName) < 1 {
		return fmt.Errorf("Must set either ChannelID or ChannelName for probe")
	}

	if len(bp.Config.ChannelID) < 1 {
		err := bp.getChannelID(bp.Config.ChannelName)
		if err != nil {
			return fmt.Errorf("could not get channel id - %s", err)
		}
	}
	if err := bp.ensureMembership(bp.Listener); err != nil {
		return fmt.Errorf("listener could not join channel - %s", err)
	}
	if err := bp.ensureMembership(bp.Speaker); err != nil {
		return fmt.Errorf("speaker could not join channel - %s", err)
	}

	bp.Listener.AddSubscription(bp)

	return nil
}

// Start will kick off the probe
func (bp *BroadcastProbe) Start() error {

	if bp.Active {
		return nil
	}

	go bp.listenForEvents()

	if bp.Config.Frequency < 0.2 {
		bp.Speaker.LogInfo("Frequency cannot be set below 0.2, setting to default 1 sec")
		bp.Config.Frequency = 1
	}

	writeTicker := time.NewTicker(time.Duration(bp.Config.Frequency * float64(time.Second)))
	go func() {
		for {
			select {
			case <-bp.StopChannel:
				bp.Active = false
				return
			case <-writeTicker.C:
				go bp.SendWrite()
			}
		}
	}()

	if bp.Config.Cutoff > 0 {
		overdueTicker := time.NewTicker(time.Duration((bp.Config.Cutoff / 4) * float64(time.Second)))
		go func() {
			for {
				select {
				case <-bp.StopChannel:
					return
				case <-overdueTicker.C:
					go bp.CheckOverdue()
				}
			}
		}()
	}

	bp.Active = true
	return nil
}

// SendWrite sends a sample post
func (bp *BroadcastProbe) SendWrite() {
	p := &model.Post{}
	uid := uuid.NewV4().String()
	sentAt := time.Now()
	bp.Messages.Add(uid, sentAt)
	p.ChannelId = bp.Config.ChannelID
	p.UserId = bp.Speaker.User.Id
	p.Message = uid
	if err := bp.Speaker.CreatePost(p); err != nil {
		bp.Speaker.LogError("Broadcast Speaker Error:", err.Error())
	}
}

func (bp *BroadcastProbe) listenForEvents() {
	for {
		select {
		case e := <-bp.EventChannel:
			bp.handleEvent(e)
		}
	}
}

func (bp *BroadcastProbe) handleEvent(event *model.WebSocketEvent) {
	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	uid := post.Message
	end := time.Now()
	start, ok := bp.Messages.Delete(uid)
	if !ok {
		bp.Speaker.Log.WithFields(logrus.Fields{
			"GUID":  uid,
			"Route": metrics.RouteBroadcastReceived,
		}).Warn("failed to find uid in probe message map")
		return
	}

	duration := end.Sub(start).Seconds()
	// Sanity protection for times, possible bad result from Message.Delete
	if duration >= bp.Config.Cutoff*2 || duration < 0 {
		bp.Messages.RLock()
		bp.Speaker.Log.WithFields(logrus.Fields{
			"GUID":     uid,
			"Route":    metrics.RouteBroadcastReceived,
			"messages": bp.Messages.Items,
			"start":    start,
			"end":      end,
			"ok":       ok,
			"post":     post,
		}).Warn("consistenty error - message length twice as long as cutoff")
		bp.Messages.RUnlock()
		return
	}

	if bp.ReportChannel == nil {
		bp.Speaker.Log.Warn("No action from BroadcastProbe handleEvent")
		return
	}

	bp.ReportChannel <- metrics.Report{
		Route:           metrics.RouteBroadcastReceived,
		DurationSeconds: duration,
	}
}

func (bp *BroadcastProbe) getChannelID(name string) error {
	channel, err := bp.Speaker.GetChannelByName(name)
	if err != nil {
		bp.Speaker.LogError("Probe error", err.Error())
		return err
	}

	bp.Config.ChannelID = channel.Id
	return nil
}

func (bp *BroadcastProbe) ensureMembership(c *mattermost.Client) error {
	err := c.JoinChannel(bp.Config.ChannelID)
	if err != nil {
		return err
	}
	return err
}

// GetSubscription adheres to SubscriptionProbe from mattermost subpackag
func (bp BroadcastProbe) GetSubscription() *mattermost.WebSocketSubscription {
	wss := mattermost.NewWebsocketSubcription(bp.EventChannel)
	// TODO: Create append helper functions
	wss.UserIDs = append(wss.UserIDs, bp.Speaker.User.Id)
	wss.ChannelIDs = append(wss.ChannelIDs, bp.Config.ChannelID)
	wss.EventTypes = append(wss.EventTypes, model.WEBSOCKET_EVENT_POSTED)
	return wss
}

func (bp BroadcastProbe) reportOverdue() {
	if bp.ReportChannel == nil {
		return
	}
	bp.ReportChannel <- metrics.Report{
		Route:           metrics.RouteBroadcastReceived,
		DurationSeconds: 0,
		Error:           fmt.Errorf("Message over cutoff %v", bp.Config.Cutoff),
	}
}

// CheckOverdue will handle any overdue messages
func (bp *BroadcastProbe) CheckOverdue() {
	if overdue := bp.Messages.Overdue(bp.Config.Cutoff); len(overdue) > 0 {
		for id := range overdue {
			bp.Messages.Delete(id)
			bp.reportOverdue()
		}
	}
}

func (bp *BroadcastProbe) String() string {
	return bp.Name
}
