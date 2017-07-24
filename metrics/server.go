package metrics

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/util"
	"github.com/prometheus/client_golang/prometheus"
)

// Server respresents the Prometheus metrics and incoming channel
type Server struct {
	ReportChannel chan TimingReport
	Log           *logrus.Logger
	Output        *logrus.Logger
}

// NewServer returns a new metric server that is ready
func NewServer(log *logrus.Logger, outputLocation string) *Server {
	tr := make(chan TimingReport, 100)

	// TODO: Move Main Logger for debugger
	var output *logrus.Logger
	if len(outputLocation) > 0 {
		output = util.NewFileLogger(outputLocation, false)
	}

	s := &Server{
		tr,
		log,
		output,
	}
	go s.MonitorTimingReports()
	return s
}

// Listen starts the prometheus server
func (s *Server) Listen(address string, port int) {
	if port == 0 {
		port = 8067
	}
	if address == "" {
		address = "0.0.0.0"
	}
	serverAddr := fmt.Sprintf("%s:%d", address, port)
	s.LogInfo("Server Started", serverAddr)

	for _, m := range ResponseHistograms {
		prometheus.MustRegister(m)
	}
	for _, m := range ResponseGauges {
		prometheus.MustRegister(m)
	}
	for _, m := range ErrorCounters {
		prometheus.MustRegister(m)
	}

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(serverAddr, nil)
}

//MonitorTimingReports loops through report channel andpasses the reports to HandleReport
func (s *Server) MonitorTimingReports() {
	for {
		select {
		case r := <-s.ReportChannel:
			s.HandleReport(r)
		}
	}
}

// HandleReport will process incoming timing reports
func (s *Server) HandleReport(r TimingReport) {
	if r.Error != nil {
		s.HandleError(r)
	}
	if _, ok := r.EnsureName(); !ok {
		s.LogWarn("HandleReport - Failed to find metric by path %v\n", r.Path)
		return
	}
	gauge, ok := ResponseGauges[r.MetricName]
	if !ok {
		s.LogWarn("HandleReport - Failed to find gauage by name %v\n", r.MetricName)
		return
	}
	histogram, ok := ResponseHistograms[r.MetricName]
	if !ok {
		s.LogWarn("HandleReport - Failed to find histogram by name %v\n", r.MetricName)
		return
	}
	gauge.Set(r.DurationSeconds)
	histogram.Observe(r.DurationSeconds)
	if s.Output != nil {
		s.Output.WithFields(logrus.Fields{
			"Metric":   r.MetricName,
			"Duration": r.DurationSeconds,
		}).Info("metric")
	}
}

// HandleError will process incoming timing reports with error
func (s *Server) HandleError(r TimingReport) {

	if _, ok := r.EnsureName(); !ok {
		s.LogWarn("HandleReport - Failed to find metric by path %v\n", r.Path)
		return
	}
	counter, ok := ErrorCounters[r.MetricName]
	if !ok {
		s.LogWarn("HandleReport - Failed to find error metric by name %v\n", r.MetricName)
		return
	}
	counter.Inc()

	if s.Output != nil {
		s.Output.WithFields(logrus.Fields{
			"Metric": r.MetricName,
			"Err":    r.Error,
		}).Info("metric error")
	}
}

// LogInfo is a shortcut for logging if the log exists
func (s *Server) LogInfo(items ...interface{}) {
	if s.Log != nil {
		s.Log.Info(items)
	}
}

// LogWarn is a shortcut for logging if the log exists
func (s *Server) LogWarn(template string, items ...interface{}) {
	if s.Log != nil {
		s.Log.Warnf(template, items)
	}
}

// LogDebug is a shortcut for logging if the log exists
func (s *Server) LogDebug(template string, items ...interface{}) {
	if s.Log != nil {
		s.Log.Infof(template, items)
	}
}
