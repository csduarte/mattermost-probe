package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// Server respresents the Prometheus metrics and incoming channel
type Server struct {
	ReportChannel TimingChannel
}

// NewServer returns a new metric server that is ready
func NewServer() *Server {
	tr := make(chan TimingReport)
	s := &Server{
		tr,
	}
	go s.MonitorTimingReports()
	return s
}

// Listen starts the prometheus server
func (s *Server) Listen() {
	for _, m := range Metrics {
		prometheus.MustRegister(m)
	}
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
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
	if len(r.MetricName) == 0 {
		mn, ok := LookupMetricNameByPath(r.Path)
		if !ok {
			fmt.Printf("WARNING: Failed to find metric by path %v\n", r.Path)
			return
		}
		r.MetricName = mn
	}
	metric, ok := Metrics[r.MetricName]
	if !ok {
		fmt.Printf("WARNING: Failed to find metric by name %v\n", r.MetricName)
		return
	}
	fmt.Println(r.MetricName, ": ", r.DurationSeconds)
	metric.Set(r.DurationSeconds)
}
